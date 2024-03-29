package main

import (
	"github.com/circonus-labs/terraform-provider-circonus/circonus"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider { //nolint
			return circonus.Provider()
		},
	})
}
