package main

import (
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Document struct {
	Lines          []string
	IRCache        *ir.IR
	Tokens         *tokens.Tokenizer
	SemanticTokens *protocol.SemanticTokens
	Config         *Config
}

var Documents = map[string]*Document{}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	doc := &Document{
		Lines: strings.Split(params.TextDocument.Text, "\n"),
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
	doc.Lines = strings.Split(c.Text, "\n")
	tokenizeDoc(doc, params.TextDocument.URI, c.Text)
	return nil
}
