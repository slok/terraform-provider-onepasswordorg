package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

var (
	_ resource.Resource                = &groupResource{}
	_ resource.ResourceWithConfigure   = &groupResource{}
	_ resource.ResourceWithImportState = &groupResource{}
)

func NewGroupResource() resource.Resource {
	return &groupResource{}
}

type groupResource struct {
	repo storage.Repository
}

func (r *groupResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group"
}

func (r *groupResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a Group resource.

A 1password group is like a team that can contain people and can be used to give access to vaults as a
group of users.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Required:    true,
				Description: "The name of the group.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Managed by Terraform"),
				Description: "The description of the group.",
			},
		},
	}
}

func (r *groupResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r *groupResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var tfGroup Group
	diags := req.Plan.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create group.
	g := mapTfToModelGroup(tfGroup)
	newGroup, err := r.repo.CreateGroup(ctx, g)
	if err != nil {
		resp.Diagnostics.AddError("Error creating group", "Could not create group, unexpected error: "+err.Error())
		return
	}

	// Map group to tf model.
	newTfGroup := mapModelToTfGroup(*newGroup)

	diags = resp.State.Set(ctx, newTfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan.
	var tfGroup Group
	diags := req.State.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group.
	id := tfGroup.ID.ValueString()
	group, err := r.repo.GetGroupByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading group", fmt.Sprintf("Could not get group %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map group to tf model.
	readTfGroup := mapModelToTfGroup(*group)

	diags = resp.State.Set(ctx, readTfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) { // Get plan values.
	var plan Group
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state.
	var state Group
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan group as the new data and set ID from state.
	g := mapTfToModelGroup(plan)
	g.ID = state.ID.ValueString()

	newGroup, err := r.repo.EnsureGroup(ctx, g)
	if err != nil {
		resp.Diagnostics.AddError("Error updating group", "Could not update group, unexpected error: "+err.Error())
		return
	}

	// Map group to tf model.
	readTfGroup := mapModelToTfGroup(*newGroup)

	diags = resp.State.Set(ctx, readTfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *groupResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan.
	var tfGroup Group
	diags := req.State.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group.
	id := tfGroup.ID.ValueString()
	err := r.repo.DeleteGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", fmt.Sprintf("Could not delete group %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r *groupResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
