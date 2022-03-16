package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

type dataSourceVaultType struct{}

func (d dataSourceVaultType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides information about a 1password vault.
`,
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "The name of the vault.",
				Required:    true,
				Validators:  []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Type:        types.StringType,
			},
			"description": {
				Computed: true,
				Type:     types.StringType,
			},
			"id": {
				Computed: true,
				Type:     types.StringType,
			},
		},
	}, nil
}

func (d dataSourceVaultType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceVault{
		p: *prv,
	}, nil
}

type dataSourceVault struct {
	p provider
}

func (d dataSourceVault) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfVault Vault
	diags := req.Config.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource.
	vault, err := d.p.repo.GetVaultByName(ctx, tfVault.Name.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error getting vault", "Could not get vault, unexpected error: "+err.Error())
		return
	}

	newTfVault := mapModelToTfVault(*vault)

	diags = resp.State.Set(ctx, newTfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
