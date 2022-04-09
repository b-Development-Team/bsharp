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

func matchPrev(tokID int, doc *Document, matchers []TokenMatcher) bool {
	if tokID < len(matchers) {
		return false
	}

	for i := 0; i < len(matchers); i++ {
		if !matchers[i].Match(doc.Tokens.Tokens[tokID-(i+1)]) {
			return false
		}
	}

	return true
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
						va := doc.IRCache.Variables[v]
						return &protocol.Hover{
							Contents: protocol.MarkupContent{
								Kind:  protocol.MarkupKindMarkdown,
								Value: fmt.Sprintf("```bsharp\n[VAR %s %s]\n```", va.Name, va.Type.String()),
							},
							Range: &protocol.Range{
								Start: protocol.Position{Line: uint32(tok.Pos.Line), Character: uint32(tok.Pos.Char)},
								End:   protocol.Position{Line: uint32(tok.Pos.EndLine), Character: uint32(tok.Pos.EndChar)},
							},
						}, nil
					}
				}
			}
		}
		return nil, nil
	}

	// Function, check if exists and find pos
	for _, fn := range getFns(doc) {
		if fn.Name == *name {
			// Gen text
			text := &strings.Builder{}
			fmt.Fprintf(text, "[FUNC %s", fn.Name)
			for i, par := range fn.Params {
				if fn.ParNames != nil {
					fmt.Fprintf(text, " [PARAM %s %s]", fn.ParNames[i], par.String())
				} else {
					fmt.Fprintf(text, " [PARAM %s]", par.String())
				}
			}
			if fn.RetType != nil && !types.NULL.Equal(fn.RetType) {
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
