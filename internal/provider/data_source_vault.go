package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

var (
	_ datasource.DataSource              = &vaultDataSource{}
	_ datasource.DataSourceWithConfigure = &vaultDataSource{}
)

func NewVaultDataSource() datasource.DataSource {
	return &vaultDataSource{}
}

type vaultDataSource struct {
	repo storage.Repository
}

func (d *vaultDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

func (d *vaultDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides information about a 1password vault.
`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the vault.",
				Required:    true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
			},
			"description": schema.StringAttribute{
				Description: "The description of the vault.",
				Computed:    true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *vaultDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	appServices := getAppServicesFromDatasourceRequest(&req)
	if appServices == nil {
		return
	}

	d.repo = appServices.Repository
}

func (d *vaultDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve values.
	var tfVault Vault
	diags := req.Config.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource.
	vault, err := d.repo.GetVaultByName(ctx, tfVault.Name.ValueString())
	if err != nil {
		resp.Diagnostics.AddError("Error getting vault", "Could not get vault, unexpected error: "+err.Error())
		return
	}

	newTfVault := mapModelToTfVault(*vault)

	diags = resp.State.Set(ctx, newTfVault)
	resp.Diagnostics.Append(diags...)
}

func getAppServicesFromDatasourceRequest(req *datasource.ConfigureRequest) *providerAppServices {
	if req.ProviderData != nil {
		if c, ok := req.ProviderData.(providerAppServices); ok {
			return &c
		}
	}

	return nil
}
