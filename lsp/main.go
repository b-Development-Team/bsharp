package main

import (
	"log"

	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
	"github.com/tliron/glsp/server"
	"github.com/tliron/kutil/logging"

	// Must include a backend implementation. See kutil's logging/ for other options.
	_ "github.com/tliron/kutil/logging/simple"
)

const lsName = "my language"

var version string = "0.0.1"
var handler protocol.Handler

func Ptr[T any](v T) *T {
	return &v
}

func main() {
	// This increases logging verbosity (optional)
	logging.Configure(1, nil)

	handler = protocol.Handler{
		Initialize:             initialize,
		Initialized:            initialized,
		Shutdown:               shutdown,
		SetTrace:               setTrace,
		TextDocumentCompletion: textDocumentCompletion,
		CompletionItemResolve:  completionItemResolve,
	}

	server := server.NewServer(&handler, lsName, false)

	server.RunStdio()
}

func textDocumentCompletion(context *glsp.Context, params *protocol.CompletionParams) (interface{}, error) {
	var completionItems []protocol.CompletionItem
	completionItems = append(completionItems, protocol.CompletionItem{
		Label: "server",
		Kind:  Ptr(protocol.CompletionItemKindFile),
	})
	log.Println("COMPLETION")
	return completionItems, nil
}

func completionItemResolve(context *glsp.Context, params *protocol.CompletionItem) (*protocol.CompletionItem, error) {
	return params, nil
}

func initialize(context *glsp.Context, params *protocol.InitializeParams) (interface{}, error) {
	capabilities := handler.CreateServerCapabilities()
	capabilities.CompletionProvider = &protocol.CompletionOptions{
		ResolveProvider:   Ptr(true),
		TriggerCharacters: []string{"#"},
	}

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
