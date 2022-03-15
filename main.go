package main

import (
	"context"
	"fmt"
	"os"

	"github.com/hashicorp/terraform-plugin-framework/tfsdk"

	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

func run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	err := tfsdk.Serve(ctx, provider.New, tfsdk.ServeOpts{
		Name: "onepasswordorg",
	})

	return err
}

func main() {
	err := run(context.Background())

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error running Terraform provider: %s", err)
		os.Exit(1)
	}
}
