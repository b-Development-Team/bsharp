package main

import (
	"path/filepath"
	"strings"

	"github.com/Nv7-Github/bsharp/bot"
	"github.com/Nv7-Github/bsharp/ir"
	"github.com/Nv7-Github/bsharp/tokens"
)

func tokenizeDoc(doc *Document, uri, text string) {
	// Tokenize
	tok := tokens.NewTokenizer(tokens.NewStream(strings.TrimPrefix(uri, RootURI), text))
	err := tok.Tokenize()
	if err != nil {
		return
	}
	doc.Tokens = tok
	updateDocTokenCache(doc)
}

func getBld(doc *Document) *ir.Builder {
	bld := ir.NewBuilder()
	if doc.Config.DiscordSupport {
		for _, ext := range bot.Exts {
			bld.AddExtension(ext)
		}
	}
	return bld
}

func buildDoc(doc *Document, uri string) {
	path := strings.TrimPrefix(uri, "file://")
	fs := &FS{}
	p, err := fs.Parse(filepath.Base(path))
	if err == nil {
		bld := getBld(doc)
		err = bld.Build(p, fs)
		if err == nil { // No error, save IR cache
			doc.IRCache = bld.IR()
		}
	}
}
