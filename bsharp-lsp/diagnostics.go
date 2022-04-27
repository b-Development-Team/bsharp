package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var Root string
var RootURI string

type FS struct{}

func (d *FS) Parse(name string) (*parser.Parser, error) {
	var src string
	doc, exists := Documents[filepath.Join(RootURI, name)]
	if exists {
		src = doc.Source
	} else {
		var err error
		dat, err := os.ReadFile(filepath.Join(Root, name))
		if err != nil {
			return nil, err
		}
		src = string(dat)
	}
	stream := tokens.NewTokenizer(tokens.NewStream(name, src))
	err := stream.Tokenize()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(stream)

	err = parser.Parse()
	return parser, err
}

func textDocumentDidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	// Run when doc is saved
	doc := Documents[params.TextDocument.URI]
	path := strings.TrimPrefix(params.TextDocument.URI, RootURI+"/")
	fs := &FS{}
	p, err := fs.Parse(path)
	if err != nil {
		return nil // Invalid code, don't run diagnostics
	}
	ir := getBld(doc)
	err = ir.Build(p, fs)
	if err == nil { // No error, save IR cache, clear diagnostics
		doc.IRCache = ir.IR()

		// Clear diagnostics
		context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
			URI:         params.TextDocument.URI,
			Diagnostics: []protocol.Diagnostic{},
		})
		return nil
	}
	v, ok := err.(*tokens.PosError)
	if !ok {
		return nil // No diagnostic error
	}

	// Make diagnostic
	context.Notify(protocol.ServerTextDocumentPublishDiagnostics, protocol.PublishDiagnosticsParams{
		URI: filepath.Join(RootURI, v.Pos.File),
		Diagnostics: []protocol.Diagnostic{
			{
				Range: protocol.Range{
					Start: protocol.Position{Line: uint32(v.Pos.Line), Character: uint32(v.Pos.Char)},
					End:   protocol.Position{Line: uint32(v.Pos.EndLine), Character: uint32(v.Pos.EndChar)},
				},
				Severity: Ptr(protocol.DiagnosticSeverityError),
				Source:   Ptr("bsharp"),
				Message:  v.Msg,
			},
		},
	})

	return nil
}
