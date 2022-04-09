package main

import (
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Document struct {
	Source         string
	IRCache        *ir.IR
	Tokens         *tokens.Tokenizer
	SemanticTokens *protocol.SemanticTokens
	Config         *Config
}

var Documents = map[string]*Document{}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	doc := &Document{
		Source: params.TextDocument.Text,
	}
	Documents[params.TextDocument.URI] = doc

	go func() {
		doc.Config = config(context, params.TextDocument.URI)

		// Build IR Cache if possible
		buildDoc(doc, params.TextDocument.URI)

		// Semantic tokens
		tokenizeDoc(doc, params.TextDocument.URI, params.TextDocument.Text)

	}()

	return nil
}

func textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	delete(Documents, params.TextDocument.URI)
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	doc := Documents[params.TextDocument.URI]
	c := params.ContentChanges[0].(protocol.TextDocumentContentChangeEventWhole)
	doc.Source = c.Text
	tokenizeDoc(doc, params.TextDocument.URI, c.Text)
	buildDoc(doc, params.TextDocument.URI)
	return nil
}
