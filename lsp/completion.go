package main

import (
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
	if tokid < 1 {
		return nil, nil
	}

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
	if !isFn {
		if matchPrev(tokid, doc, []TokenMatcher{TokMatchTV(tokens.TokenTypeIdent, "VAR"), TokMatchT(tokens.TokenTypeLBrack)}) {
			// Variable
		}
	}

	// Get items
	out := make([]protocol.CompletionItem, 0)
	for _, fn := range fns {
		if strings.HasPrefix(fn.Name, tok.Value) {
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
