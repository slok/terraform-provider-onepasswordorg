package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type dataSourceUserType struct{}

func (d dataSourceUserType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides information about a 1password user.
`,
		Attributes: map[string]tfsdk.Attribute{
			"email": {
				Description: "The email of the user.",
				Required:    true,
				Type:        types.StringType,
			},
			"name": {
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

func (d dataSourceUserType) NewDataSource(ctx context.Context, p tfsdk.Provider) (tfsdk.DataSource, diag.Diagnostics) {
	prv := p.(*provider)
	return dataSourceUser{
		p: *prv,
	}, nil
}

type dataSourceUser struct {
	p provider
}

func (d dataSourceUser) Read(ctx context.Context, req tfsdk.ReadDataSourceRequest, resp *tfsdk.ReadDataSourceResponse) {
	if !d.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values.
	var tfUser User
	diags := req.Config.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	email := tfUser.Email.Value
	if email == "" {
		resp.Diagnostics.AddError("Ivalid email", "Could not get user, because email is empty")
		return
	}

	// Get user.
	user, err := d.p.repo.GetUserByEmail(ctx, tfUser.Email.Value)
	if err != nil {
		resp.Diagnostics.AddError("Error getting user", "Could not get user, unexpected error: "+err.Error())
		return
	}

	newTfUser := mapModelToTfUser(*user)

	diags = resp.State.Set(ctx, newTfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}
