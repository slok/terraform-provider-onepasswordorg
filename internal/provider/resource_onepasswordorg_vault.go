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
	_ resource.Resource                = &vaultResource{}
	_ resource.ResourceWithConfigure   = &vaultResource{}
	_ resource.ResourceWithImportState = &vaultResource{}
)

func NewVaultResource() resource.Resource {
	return &vaultResource{}
}

type vaultResource struct {
	repo storage.Repository
}

func (r *vaultResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault"
}

func (r *vaultResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a vault resource.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "The name of the vault.",
			},
			"description": schema.StringAttribute{
				Optional:    true,
				Computed:    true,
				Default:     stringdefault.StaticString("Managed by Terraform"),
				Description: "The description of the vault.",
			},
		},
	}
}

func (r *vaultResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r *vaultResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var tfVault Vault
	diags := req.Plan.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create vault.
	v := mapTfToModelVault(tfVault)
	newVault, err := r.repo.CreateVault(ctx, v)
	if err != nil {
		resp.Diagnostics.AddError("Error creating vault", "Could not create vault, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfVault := mapModelToTfVault(*newVault)

	diags = resp.State.Set(ctx, newTfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vaultResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan.
	var tfVault Vault
	diags := req.State.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource.
	id := tfVault.ID.ValueString()
	vault, err := r.repo.GetVaultByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading vault", fmt.Sprintf("Could not get vault %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map resource to tf model.
	readTfVault := mapModelToTfVault(*vault)

	diags = resp.State.Set(ctx, readTfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vaultResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Get plan values.
	var plan Vault
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state.
	var state Vault
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan group as the new data and set ID from state.
	v := mapTfToModelVault(plan)
	v.ID = state.ID.ValueString()

	newVault, err := r.repo.EnsureVault(ctx, v)
	if err != nil {
		resp.Diagnostics.AddError("Error updating vault", "Could not update vault, unexpected error: "+err.Error())
		return
	}

	// Map vault to tf model.
	readTfVault := mapModelToTfVault(*newVault)

	diags = resp.State.Set(ctx, readTfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vaultResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan.
	var tfVault Vault
	diags := req.State.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete resource.
	id := tfVault.ID.ValueString()
	err := r.repo.DeleteVault(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting vault", fmt.Sprintf("Could not delete vault %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r *vaultResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func getAppServicesFromResourceRequest(req *resource.ConfigureRequest) *providerAppServices {
	if req.ProviderData != nil {
		if c, ok := req.ProviderData.(providerAppServices); ok {
			return &c
		}
	}

	return nil
}
