package provider

import (
	"context"

	"github.com/hashicorp/terraform-plugin-framework/datasource"
	"github.com/hashicorp/terraform-plugin-framework/datasource/schema"

	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

var (
	_ datasource.DataSource              = &groupDataSource{}
	_ datasource.DataSourceWithConfigure = &groupDataSource{}
)

func NewGroupDataSource() datasource.DataSource {
	return &groupDataSource{}
}

type groupDataSource struct {
	repo storage.Repository
}

func (d *groupDataSource) Metadata(_ context.Context, req datasource.MetadataRequest, resp *datasource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (d *groupDataSource) Schema(_ context.Context, _ datasource.SchemaRequest, resp *datasource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides information about a 1password group.
`,
		Attributes: map[string]schema.Attribute{
			"name": schema.StringAttribute{
				Description: "The name of the group.",
				Required:    true,
			},
			"description": schema.StringAttribute{
				Computed: true,
			},
			"id": schema.StringAttribute{
				Computed: true,
			},
		},
	}
}

func (d *groupDataSource) Configure(ctx context.Context, req datasource.ConfigureRequest, resp *datasource.ConfigureResponse) {
	appServices := getAppServicesFromDatasourceRequest(&req)
	if appServices == nil {
		return
	}

	d.repo = appServices.Repository
}

func (d *groupDataSource) Read(ctx context.Context, req datasource.ReadRequest, resp *datasource.ReadResponse) {
	// Retrieve values.
	var tfGroup Group
	diags := req.Config.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group.
	group, err := d.repo.GetGroupByName(ctx, tfGroup.Name.ValueString())
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
