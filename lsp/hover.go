package main

import (
	"fmt"
	"strings"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/Nv7-Github/bsharp/types"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func docHover(context *glsp.Context, params *protocol.HoverParams) (*protocol.Hover, error) {
	// Find token
	doc, exists := Documents[params.TextDocument.URI]
	if !exists {
		return nil, nil
	}
	if doc.Tokens == nil || doc.IRCache == nil {
		return nil, nil
	}
	id := -1
	for i, tok := range doc.Tokens.Tokens {
		if tok.Pos.Line > int(params.Position.Line) || (tok.Pos.Line == int(params.Position.Line) && tok.Pos.Char > int(params.Position.Character)) {
			break
		}
		id = i
	}
	if id < 1 {
		return nil, nil
	}

	// Check if function
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
