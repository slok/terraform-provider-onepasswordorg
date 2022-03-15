package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceGroupType struct{}

func (d dataSourceGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides information about a 1password group.
`,
		Attributes: map[string]tfsdk.Attribute{
			"name": {
				Description: "The name of the group.",
				Required:    true,
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

func (d dataSourceGroupType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceGroup{
		p: *prv,
	}, nil
}

type dataSourceGroup struct {
	p provider
}

func (d dataSourceGroup) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfGroup Group
	diags := req.Config.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	name := tfGroup.Name.Value
	if name == "" {
		resp.Diagnostics.AddError("Ivalid name", "Could not get group, because name is empty")
		return
	}

	// Get group.
	group, err := d.p.repo.GetGroupByName(ctx, name)
	if err != nil {
		resp.Diagnostics.AddError("Error getting group", "Could not get group, unexpected error: "+err.Error())
		return
	}

	newTfGroup := mapModelToTfGroup(*group)

	diags = resp.State.Set(ctx, newTfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
