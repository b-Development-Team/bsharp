package main

import (
	"github.com/tliron/glsp"
	protocol "github.com/tliron/glsp/protocol_3_16"
)

type Config struct {
	DiscordSupport bool
}

func config(context *glsp.Context, doc string) *Config {
	var out []bool
	context.Call(protocol.ServerWorkspaceConfiguration, &protocol.ConfigurationParams{
		Items: []protocol.ConfigurationItem{
			{
				ScopeURI: &doc,
				Section:  Ptr("discordSupport"),
			},
		},
	}, &out)
	return &Config{
		DiscordSupport: out[0],
	}
}
