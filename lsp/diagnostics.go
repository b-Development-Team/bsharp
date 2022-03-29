package main

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/parser"
	"github.com/Nv7-Github/bsharp/tokens"
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

var Root string

type FS struct{}

func (d *FS) Parse(name string) (*parser.Parser, error) {
	src, err := os.ReadFile(filepath.Join(Root, name))
	if err != nil {
		return nil, err
	}
	stream := tokens.NewTokenizer(tokens.NewStream(name, string(src)))
	err = stream.Tokenize()
	if err != nil {
		return nil, err
	}
	parser := parser.NewParser(stream)

	err = parser.Parse()
	return parser, err
}

func textDocumentDidSave(context *glsp.Context, params *protocol.DidSaveTextDocumentParams) error {
	// Run when doc is saved
	path := strings.TrimPrefix(params.TextDocument.URI, "file://")
	fs := &FS{}
	p, err := fs.Parse(filepath.Base(path))
	if err != nil {
		return nil // Invalid code, don't run diagnostics
	}
	ir := ir.NewBuilder()
	err = ir.Build(p, fs)
	if err == nil { // No error, save IR cache, clear diagnostics
		doc := Documents[params.TextDocument.URI]
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
		URI: params.TextDocument.URI,
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
