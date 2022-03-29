package main

import (
	"path/filepath"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Document struct {
	Lines   []string
	IRCache *ir.IR
}

var Documents = map[string]*Document{}

func textDocumentDidOpen(context *glsp.Context, params *protocol.DidOpenTextDocumentParams) error {
	doc := &Document{
		Lines: strings.Split(params.TextDocument.Text, "\n"),
	}
	Documents[params.TextDocument.URI] = doc

	// Build IR Cache if possible
	path := strings.TrimPrefix(params.TextDocument.URI, "file://")
	fs := &FS{}
	p, err := fs.Parse(filepath.Base(path))
	if err == nil {
		bld := ir.NewBuilder()
		err = bld.Build(p, fs)
		if err == nil { // No error, save IR cache
			doc.IRCache = bld.IR()
		}
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
		c := change.(protocol.TextDocumentContentChangeEventWhole)
		doc.Lines = strings.Split(c.Text, "\n")
	}
	return nil
}
