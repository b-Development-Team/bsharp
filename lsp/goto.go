package main

import (
	"path/filepath"

	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

func definitionFn(context *glsp.Context, params *protocol.DefinitionParams) (interface{}, error) {
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
	tok := doc.Tokens.Tokens[id]
	if matchPrev(id, doc, []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack)}) {
		name = &tok.Value
	} else if matchPrev(id, doc, []TokenMatcher{TokMatchTV(tokens.TokenTypeIdent, "FN"), TokMatchT(tokens.TokenTypeLBrack)}) {
		name = &tok.Value
	}
	if name == nil {
		// Variable?
		if matchPrev(id, doc, []TokenMatcher{TokMatchTV(tokens.TokenTypeIdent, "VAR"), TokMatchT(tokens.TokenTypeLBrack)}) {
			// Variable
			scope := GetScope(doc.IRCache, tok.Pos)
			if scope == nil {
				return nil, nil
			}
			for i := len(scope.Frames) - 1; i >= 0; i-- {
				for _, v := range scope.Frames[i].Variables {
					if doc.IRCache.Variables[v].Name == tok.Value {
						// Get var pos
						va := doc.IRCache.Variables[v]
						id := tokID(protocol.Position{
							Line:      uint32(va.Pos.Line),
							Character: uint32(va.Pos.Char),
						}, doc)
						name := doc.Tokens.Tokens[id+2]

						return protocol.Location{
							URI: filepath.Join(RootURI, name.Pos.File),
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
			}
		}
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
