package main

import (
	"log"
	"strings"

	"github.com/Nv7-Github/bsharp/bot"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var fns = ir.BuiltinFns()

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
		log.Println(scope, word)
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
	for _, fn := range fns {
		if strings.HasPrefix(fn.Name, word) {
			out = append(out, protocol.CompletionItem{
				Label: fn.Name,
				Kind:  Ptr(protocol.CompletionItemKindFunction),
			})
		}
	}
	if doc.Config.DiscordSupport {
		for _, fn := range bot.Exts {
			if strings.HasPrefix(fn.Name, word) {
				out = append(out, protocol.CompletionItem{
					Label: fn.Name,
					Kind:  Ptr(protocol.CompletionItemKindFunction),
				})
			}
		}
	}
	if doc.IRCache != nil {
		for _, fn := range doc.IRCache.Funcs {
			if strings.HasPrefix(fn.Name, word) {
				out = append(out, protocol.CompletionItem{
					Label: fn.Name,
					Kind:  Ptr(protocol.CompletionItemKindFunction),
				})
			}
		}
	}

	return out, nil
}
