package main

import (
	"github.com/circonus-labs/terraform-provider-circonus/circonus"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: circonus.Provider})
}
