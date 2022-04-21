package main

import (
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type TokenMatcher interface {
	Match(tokens.Token) bool
}

type TokenMatcherType struct {
	typ tokens.TokenType
}

func TokMatchT(typ tokens.TokenType) TokenMatcherType {
	return TokenMatcherType{typ}
}

func (t TokenMatcherType) Match(v tokens.Token) bool {
	return v.Typ == t.typ
}

type TokenMatcherTypeVal struct {
	typ tokens.TokenType
	val string
}

func TokMatchTV(typ tokens.TokenType, val string) TokenMatcherTypeVal {
	return TokenMatcherTypeVal{typ, val}
}

func (t TokenMatcherTypeVal) Match(v tokens.Token) bool {
	return v.Typ == t.typ && v.Value == t.val
}

const (
	MatchedTokenNone     protocol.UInteger = 3
	MatchedTokenType     protocol.UInteger = 0
	MatchedTokenVariable protocol.UInteger = 1
	MatchedTokenFunc     protocol.UInteger = 2
)

type Matcher struct {
	Matchers []TokenMatcher
	Results  []protocol.UInteger
}

var matchers = map[string][]Matcher{
	"[": {
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "VAR"), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenVariable},
		},
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "DEFINE"), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenVariable},
		},
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "FN"), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenFunc},
		},
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "MAKE"), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenType},
		},
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "PARAM"), TokMatchT(tokens.TokenTypeIdent), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenVariable, MatchedTokenType},
		},
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "RETURNS"), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenType},
		},
		{
			Matchers: []TokenMatcher{TokMatchT(tokens.TokenTypeLBrack), TokMatchTV(tokens.TokenTypeIdent, "TYPEDEF"), TokMatchT(tokens.TokenTypeIdent), TokMatchT(tokens.TokenTypeIdent)},
			Results:  []protocol.UInteger{MatchedTokenNone, MatchedTokenNone, MatchedTokenNone, MatchedTokenType},
		},
	},
}

func updateDocTokenCache(doc *Document) {
	toks := make([]protocol.UInteger, 0)

	line := 0
	char := 0

	for i, tok := range doc.Tokens.Tokens {
		// Apply rules
		v, exists := matchers[tok.Value]
		if exists {
			for _, m := range v {
				// Check if can match
				canMatch := false
				if (len(doc.Tokens.Tokens) - (i + 1)) >= len(m.Matchers) {
					canMatch = true
					for j := 0; j < len(m.Matchers); j++ {
						if !m.Matchers[j].Match(doc.Tokens.Tokens[j+i]) {
							canMatch = false
							break
						}
					}
				}

				if canMatch {
					for j, v := range m.Results {
						if v != MatchedTokenNone {
							tok := doc.Tokens.Tokens[i+j]
							if line != tok.Pos.Line {
								char = 0
							}
							vals := []protocol.UInteger{
								uint32(tok.Pos.Line - line),
								uint32(tok.Pos.Char - char),
								uint32(len(tok.Value)),
								v,
								0,
							}
							line = tok.Pos.Line
							char = tok.Pos.Char

							toks = append(toks, vals...)
						}
					}
				}
			}
		}
	}

	doc.SemanticTokens = &protocol.SemanticTokens{
		Data: toks,
	}
}

func semanticTokensFull(context *glsp.Context, params *protocol.SemanticTokensParams) (*protocol.SemanticTokens, error) {
	v, exists := Documents[params.TextDocument.URI]
	if !exists {
		return nil, nil
	}
	return v.SemanticTokens, nil
}
