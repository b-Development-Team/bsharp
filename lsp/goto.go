package main

import (
	"path/filepath"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func definitionFn(context *glsp.Context, params *protocol.DefinitionParams) (interface{}, error) {
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
			// Get name pos
			id := tokID(protocol.Position{
				Line:      uint32(fn.Pos().Line),
				Character: uint32(fn.Pos().Char),
			}, doc)
			name := doc.Tokens.Tokens[id+2]

			return protocol.Location{
				URI: filepath.Join(RootURI, fn.Pos().File),
				Range: protocol.Range{
					Start: protocol.Position{
						Line:      uint32(name.Pos.Line),
						Character: uint32(name.Pos.Char),
					},
					End: protocol.Position{
						Line:      uint32(name.Pos.EndLine),
						Character: uint32(name.Pos.EndChar),
					},
				},
			}, nil
		}
	}

	return nil, nil
}
