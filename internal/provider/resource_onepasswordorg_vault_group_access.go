package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

var (
	_ resource.Resource                = &vaultUserAccessResource{}
	_ resource.ResourceWithConfigure   = &vaultUserAccessResource{}
	_ resource.ResourceWithImportState = &vaultUserAccessResource{}
)

func NewVaultGroupAccessResource() resource.Resource {
	return &vaultGroupAccessResource{}
}

type vaultGroupAccessResource struct {
	repo storage.Repository
}

func (r *vaultGroupAccessResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_group_access"
}

func (r *vaultGroupAccessResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides vault access for a group.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"vault_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "The vault ID.",
			},
			"group_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "The group ID.",
			},
			"permissions": permissionsAttribute,
		},
	}
}

func (r *vaultGroupAccessResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r vaultGroupAccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var tfvga VaultGroupAccess
	diags := req.Plan.Get(ctx, &tfvga)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create access.
	v, err := mapTfToModelVaultGroupAccess(tfvga)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	err = r.repo.EnsureVaultGroupAccess(ctx, *v)
	if err != nil {
		resp.Diagnostics.AddError("Error creating access", "Could not create access, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfAccess, err := mapModelToTfVaultGroupAccess(*v)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map acceess:"+err.Error())
		return
	}

	// Set on state.
	diags = resp.State.Set(ctx, newTfAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r vaultGroupAccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan.
	var tfVaultGroupAccess VaultGroupAccess
	diags := req.State.Get(ctx, &tfVaultGroupAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get access.
	id := tfVaultGroupAccess.ID.ValueString()
	vaultID, groupID, err := unpackVaultGroupAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	access, err := r.repo.GetVaultGroupAccessByID(ctx, vaultID, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading access", fmt.Sprintf("Could not get access %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map to tf model.
	readTfVaultGroupAccess, err := mapModelToTfVaultGroupAccess(*access)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	diags = resp.State.Set(ctx, readTfVaultGroupAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r vaultGroupAccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan VaultGroupAccess
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VaultGroupAccess
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan as the new data and set ID from state.
	plan.ID = state.ID

	// Update access.
	v, err := mapTfToModelVaultGroupAccess(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	err = r.repo.EnsureVaultGroupAccess(ctx, *v)
	if err != nil {
		resp.Diagnostics.AddError("Error updating access", "Could not create access, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfAccess, err := mapModelToTfVaultGroupAccess(*v)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	// Set on state.
	diags = resp.State.Set(ctx, newTfAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r vaultGroupAccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan.
	var tfVaultGroupAccess VaultGroupAccess
	diags := req.State.Get(ctx, &tfVaultGroupAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete resource.
	id := tfVaultGroupAccess.ID.ValueString()
	vaultID, groupID, err := unpackVaultGroupAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	err = r.repo.DeleteVaultGroupAccess(ctx, vaultID, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting access", fmt.Sprintf("Could not delete access %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r *vaultGroupAccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTfToModelVaultGroupAccess(m VaultGroupAccess) (*model.VaultGroupAccess, error) {
	groupID := m.GroupID.ValueString()
	vaultID := m.VaultID.ValueString()

	// Check the ID is correct.
	if m.ID.ValueString() != "" {
		vid, gid, err := unpackVaultGroupAccessID(m.ID.ValueString())
		if err != nil {
			return nil, err
		}

		if gid != groupID {
			return nil, fmt.Errorf("resource id is wrong based on group ID")
		}

		if vid != vaultID {
			return nil, fmt.Errorf("resource id is wrong based on vault ID")
		}
	}

	return &model.VaultGroupAccess{
		VaultID:     vaultID,
		GroupID:     groupID,
		Permissions: mapTfToModelAccessPermissions(*m.Permissions),
	}, nil
}

func mapModelToTfVaultGroupAccess(m model.VaultGroupAccess) (*VaultGroupAccess, error) {
	id := packVaultGroupAccessID(m.VaultID, m.GroupID)

	return &VaultGroupAccess{
		ID:          types.StringValue(id),
		GroupID:     types.StringValue(m.GroupID),
		VaultID:     types.StringValue(m.VaultID),
		Permissions: mapModelToTfAccessPermissions(m.Permissions),
	}, nil
}

func packVaultGroupAccessID(vaultID, groupID string) string {
	return vaultID + "/" + groupID
}

func unpackVaultGroupAccessID(id string) (vaultID, groupID string, err error) {
	s := strings.SplitN(id, "/", 2)
	if len(s) != 2 {
		return "", "", fmt.Errorf(
			"invalid vault group access ID format: %s (expected <VAULT ID>/<GROUP ID>)", id)
	}

	return s[0], s[1], nil
}
