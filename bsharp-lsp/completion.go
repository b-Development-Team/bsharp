package main

import (
	"strings"

	"github.com/Nv7-Github/bsharp/bot"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Fn struct {
	Name   string
	Params []types.Type

	RetType  types.Type // optional
	ParNames []string   // optional
}

func getFns(doc *Document) []*Fn {
	builtin := ir.BuiltinFns()
	out := make([]*Fn, 0, len(doc.IRCache.Funcs)+len(builtin))
	for _, fn := range builtin {
		out = append(out, &Fn{
			Name:   fn.Name,
			Params: fn.Params,
		})
	}
	for _, fn := range doc.IRCache.Funcs {
		pars := make([]types.Type, len(fn.Params))
		names := make([]string, len(fn.Params))
		for i, par := range fn.Params {
			pars[i] = par.Type
			names[i] = par.Name
		}
		out = append(out, &Fn{
			Name:     fn.Name,
			Params:   pars,
			RetType:  fn.RetType,
			ParNames: names,
		})
	}
	if doc.Config != nil && doc.Config.DiscordSupport {
		for _, fn := range bot.Exts {
			out = append(out, &Fn{
				Name:    fn.Name,
				Params:  fn.Params,
				RetType: fn.RetType,
			})
		}
	}
	return out
}

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	doc := Documents[params.TextDocument.URI]
	tokid := tokID(params.Position, doc)
	if tokid < 2 {
		return nil, nil
	}
	tokid-- // Matches RBrack

	word := ""
	isFn := true
	tok := doc.Tokens.Tokens[tokid]
	if tok.Typ != tokens.TokenTypeLBrack {
		prev := doc.Tokens.Tokens[tokid-1]
		if prev.Typ == tokens.TokenTypeLBrack {
			word = tok.Value
		} else {
			isFn = false
		}
	}
	if isFn && doc.Tokens.Tokens[tokid+1].Typ == tokens.TokenTypeIdent {
		word = doc.Tokens.Tokens[tokid+1].Value
	}

	// Is variable?
	m := []TokenMatcher{TokMatchTV(tokens.TokenTypeIdent, "VAR"), TokMatchT(tokens.TokenTypeLBrack)}
	matchFu := matchPrev(tokid+1, doc, m)
	if matchPrev(tokid, doc, m) || matchFu {
		if matchFu {
			tokid++
		}

		// Find scope
		if doc.IRCache == nil {
			return nil, nil
		}
		scope := GetScope(doc.IRCache, tok.Pos)
		out := make([]protocol.CompletionItem, 0)
		done := make(map[int]struct{})
		word := ""
		tok := doc.Tokens.Tokens[tokid]
		if tok.Typ != tokens.TokenTypeRBrack {
			word = tok.Value // Already typed something
		}
		for _, scope := range scope.Frames {
			for _, v := range scope.Variables {
				_, exists := done[v]
				if exists {
					continue
				}
				done[v] = struct{}{}

				// Check if its ok
				va := doc.IRCache.Variables[v]
				if va.Pos.Line > tok.Pos.Line || (va.Pos.Line == tok.Pos.Line && va.Pos.Char > tok.Pos.Char) { // Defined after
					continue
				}

				// Check if its a match
				if strings.HasPrefix(va.Name, word) {
					out = append(out, protocol.CompletionItem{
						Label:  va.Name,
						Kind:   Ptr(protocol.CompletionItemKindVariable),
						Detail: Ptr(va.Type.String()),
					})
				}
			}
		}
		return out, nil
	}

	// Not variable, match function
	out := make([]protocol.CompletionItem, 0)
	fns := getFns(doc)
	for _, fn := range fns {
		if strings.HasPrefix(fn.Name, word) {
			out = append(out, protocol.CompletionItem{
				Label: fn.Name,
				Kind:  Ptr(protocol.CompletionItemKindFunction),
			})
		}
	}

	return out, nil
}
