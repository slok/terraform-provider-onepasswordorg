package main

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/plugin"

	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

const providerName = "registry.terraform.io/slok/onepasswordorg"

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: func() *schema.Provider {
			return provider.Provider()
		},
	})
}
