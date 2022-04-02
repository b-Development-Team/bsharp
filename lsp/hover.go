package main

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func tokID(pos protocol.Position, doc *Document) int {
	if doc.Tokens == nil || doc.IRCache == nil {
		return -1
	}
	id := -1
	for i, tok := range doc.Tokens.Tokens {
		if tok.Pos.Line > int(pos.Line) || (tok.Pos.Line == int(pos.Line) && tok.Pos.Char > int(pos.Character)) {
			break
		}
		id = i
	}
	return id
}

func docHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	// Find token
	doc, exists := Documents[params.TextDocument.URI]
	if !exists {
		return nil, nil
	}
	id := tokID(params.Position, doc)
	if id < 0 {
		return nil, nil
	}

	// Check if function, get name
	var name *string = nil
	prev := doc.Tokens.Tokens[id-1]
	tok := doc.Tokens.Tokens[id]
	if prev.Typ == tokens.TokenTypeLBrack {
		name = &tok.Value
	} else if prev.Typ == tokens.TokenTypeIdent && id > 1 {
		prev2 := doc.Tokens.Tokens[id-2]
		if prev2.Typ == tokens.TokenTypeLBrack {
			name = &tok.Value
		}
	}
	if name == nil {
		return nil, nil
	}

	// Function, check if exists and find pos
	for _, fn := range doc.IRCache.Funcs {
		if fn.Name == *name {
			// Gen text
			text := &strings.Builder{}
			fmt.Fprintf(text, "[FUNC %s", fn.Name)
			for _, par := range fn.Params {
				fmt.Fprintf(text, " [PARAM %s %s]", par.Name, par.Type.String())
			}
			if !types.NULL.Equal(fn.RetType) {
				fmt.Fprintf(text, " [RETURNS %s]", fn.RetType.String())
			}
			text.WriteString("]")

			return &protocol.Hover{
				Contents: protocol.MarkupContent{
					Kind:  protocol.MarkupKindMarkdown,
					Value: fmt.Sprintf("```bsharp\n%s\n```", text.String()),
				},
				Range: &protocol.Range{
					Start: protocol.Position{Line: uint32(tok.Pos.Line), Character: uint32(tok.Pos.Char)},
					End:   protocol.Position{Line: uint32(tok.Pos.EndLine), Character: uint32(tok.Pos.EndChar)},
				},
			}, nil
		}
	}

	return nil, nil
}
