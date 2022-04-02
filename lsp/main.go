package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"

	// Must include a backend implementation. See kutil's logging/ for other options.
	_ "github.com/tliron/kutil/logging/simple"
)

const lsName = "B#"

var version string = "0.0.1"
var handler protocol.Handler

func Ptr[T any](v T) *T {
	return &v
}

func main() {
	// This increases logging verbosity (optional)
	logging.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:                     initialize,
		Initialized:                    initialized,
		Shutdown:                       shutdown,
		SetTrace:                       setTrace,
		TextDocumentCompletion:         textDocumentCompletion,
		TextDocumentDidOpen:            textDocumentDidOpen,
		TextDocumentDidChange:          textDocumentDidChange,
		TextDocumentDidClose:           textDocumentDidClose,
		TextDocumentDidSave:            textDocumentDidSave,
		TextDocumentSemanticTokensFull: semanticTokensFull,
		TextDocumentHover:              docHover,
		TextDocumentDefinition:         definitionFn,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		TriggerCharacters: []string{"["},
	}
	capabilities.TextDocumentSync = Ptr(protocol.TextDocumentSyncKindFull)
	capabilities.SemanticTokensProvider = protocol.SemanticTokensOptions{
		Legend: protocol.SemanticTokensLegend{
			TokenTypes:     []string{string(protocol.SemanticTokenTypeType), string(protocol.SemanticTokenTypeVariable), string(protocol.SemanticTokenTypeFunction)},
			TokenModifiers: []string{},
		},
		Range: false,
		Full:  true,
	}
	capabilities.HoverProvider = true
	capabilities.DefinitionProvider = true

	Root = *params.RootPath
	RootURI = *params.RootURI

	return protocol.InitializeResult{
		Capabilities: capabilities,
		ServerInfo: &protocol.InitializeResultServerInfo{
			Name:    lsName,
			Version: &version,
		},
	}, nil
}

func initialized(context *glsp.Context, params *protocol.InitializedParams) error {
	return nil
}

func shutdown(context *glsp.Context) error {
	protocol.SetTraceValue(protocol.TraceValueOff)
	return nil
}

func setTrace(context *glsp.Context, params *protocol.SetTraceParams) error {
	protocol.SetTraceValue(params.Value)
	return nil
}
