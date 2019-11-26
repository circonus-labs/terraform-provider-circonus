package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
	"github.com/terraform-providers/terraform-provider-circonus/circonus"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: circonus.Provider})
}
