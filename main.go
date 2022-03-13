package main

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

func main() {
	tfsdk.Serve(context.Background(), provider.New, tfsdk.ServeOpts{
		Name: "onepasswordorg",
	})
}
