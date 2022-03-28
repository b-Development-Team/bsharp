package main

import (
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var fns = ir.BuiltinFns()

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	doc := Documents[params.TextDocument.URI]
	pos := params.Position.Character
	line := doc.Lines[params.Position.Line]

	// Find the word
	var word string
	for i := pos - 1; i >= 0; i-- {
		if line[i] == '[' {
			break
		}
		word = string(line[i]) + word
	}

	// Get items
	out := make([]protocol.CompletionItem, 0)
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
