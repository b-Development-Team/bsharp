package main

import (
	"strings"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Document struct {
	Lines []string
}

var Documents = map[string]*Document{}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	Documents[params.TextDocument.URI] = &Document{
		Lines: strings.Split(params.TextDocument.Text, "\n"),
	}
	return nil
}

func textDocumentDidClose(context *glsp.Context, params *protocol.DidCloseTextDocumentParams) error {
	delete(Documents, params.TextDocument.URI)
	return nil
}

func textDocumentDidChange(context *glsp.Context, params *protocol.DidChangeTextDocumentParams) error {
	doc := Documents[params.TextDocument.URI]
	for _, change := range params.ContentChanges {
		c := change.(protocol.TextDocumentContentChangeEvent)
		for i := c.Range.Start.Line; i <= c.Range.End.Line; i++ {
			if i < uint32(len(doc.Lines)) {
				// Update line
				if c.Range.Start.Character == 0 && c.Range.End.Character == 0 {
					// Replace whole line
					doc.Lines[i] = c.Text
				} else {
					// Replace part of line
					doc.Lines[i] = doc.Lines[i][:c.Range.Start.Character] + c.Text + doc.Lines[i][c.Range.End.Character:]
				}
			} else {
				doc.Lines = append(doc.Lines, c.Text)
			}
		}
	}
	return nil
}
