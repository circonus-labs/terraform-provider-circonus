package main

import (
	"github.com/hashicorp/terraform/plugin"
	"github.com/rileyberton/terraform-provider-circonus/circonus"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: circonus.Provider})
}
