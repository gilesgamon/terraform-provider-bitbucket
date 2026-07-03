package main

import (
	"flag"

	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
	"github.com/terraform-providers/terraform-provider-bitbucket/bitbucket"
)

// version is set at build time via -ldflags "-X main.version=...".
var version = "dev"

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	bitbucket.ProviderVersion = version

	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: bitbucket.Provider,
		ProviderAddr: "DrFaust92/bitbucket",
		Debug:        debug,
	})
}
