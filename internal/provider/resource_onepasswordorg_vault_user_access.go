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

func NewVaultUserAccessResource() resource.Resource {
	return &vaultUserAccessResource{}
}

type vaultUserAccessResource struct {
	repo storage.Repository
}

func (r *vaultUserAccessResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_vault_user_access"
}

func (r *vaultUserAccessResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides vault access for a user.
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
			"user_id": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "The user ID.",
			},
			"permissions": permissionsAttribute,
		},
	}
}

func (r *vaultUserAccessResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r *vaultUserAccessResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var tfvga VaultUserAccess
	diags := req.Plan.Get(ctx, &tfvga)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create access.
	v, err := mapTfToModelVaultUserAccess(tfvga)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	err = r.repo.EnsureVaultUserAccess(ctx, *v)
	if err != nil {
		resp.Diagnostics.AddError("Error creating access", "Could not create access, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfAccess, err := mapModelToTfVaultUserAccess(*v)
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

func (r *vaultUserAccessResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan.
	var tfVaultUserAccess VaultUserAccess
	diags := req.State.Get(ctx, &tfVaultUserAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get access.
	id := tfVaultUserAccess.ID.ValueString()
	vaultID, userID, err := unpackVaultUserAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	access, err := r.repo.GetVaultUserAccessByID(ctx, vaultID, userID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading access", fmt.Sprintf("Could not get access %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map to tf model.
	readTfVaultUserAccess, err := mapModelToTfVaultUserAccess(*access)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	diags = resp.State.Set(ctx, readTfVaultUserAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r *vaultUserAccessResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan VaultUserAccess
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state VaultUserAccess
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan as the new data and set ID from state.
	plan.ID = state.ID

	// Update access.
	v, err := mapTfToModelVaultUserAccess(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping access", "Could not map access:"+err.Error())
		return
	}

	err = r.repo.EnsureVaultUserAccess(ctx, *v)
	if err != nil {
		resp.Diagnostics.AddError("Error updating access", "Could not create access, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfAccess, err := mapModelToTfVaultUserAccess(*v)
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

func (r *vaultUserAccessResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan.
	var tfVaultUserAccess VaultUserAccess
	diags := req.State.Get(ctx, &tfVaultUserAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete resource.
	id := tfVaultUserAccess.ID.ValueString()
	vaultID, userID, err := unpackVaultUserAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	err = r.repo.DeleteVaultUserAccess(ctx, vaultID, userID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting access", fmt.Sprintf("Could not delete access %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r *vaultUserAccessResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTfToModelVaultUserAccess(m VaultUserAccess) (*model.VaultUserAccess, error) {
	userID := m.UserID.ValueString()
	vaultID := m.VaultID.ValueString()

	// Check the ID is correct.
	if m.ID.ValueString() != "" {
		vid, uid, err := unpackVaultUserAccessID(m.ID.ValueString())
		if err != nil {
			return nil, err
		}

		if uid != userID {
			return nil, fmt.Errorf("resource id is wrong based on user ID")
		}

		if vid != vaultID {
			return nil, fmt.Errorf("resource id is wrong based on vault ID")
		}
	}

	return &model.VaultUserAccess{
		VaultID:     vaultID,
		UserID:      userID,
		Permissions: mapTfToModelAccessPermissions(*m.Permissions),
	}, nil
}

func mapModelToTfVaultUserAccess(m model.VaultUserAccess) (*VaultUserAccess, error) {
	id := packVaultUserAccessID(m.VaultID, m.UserID)

	return &VaultUserAccess{
		ID:          types.StringValue(id),
		UserID:      types.StringValue(m.UserID),
		VaultID:     types.StringValue(m.VaultID),
		Permissions: mapModelToTfAccessPermissions(m.Permissions),
	}, nil
}

func packVaultUserAccessID(vaultID, userID string) string {
	return vaultID + "/" + userID
}

func unpackVaultUserAccessID(id string) (vaultID, userID string, err error) {
	s := strings.SplitN(id, "/", 2)
	if len(s) != 2 {
		return "", "", fmt.Errorf(
			"invalid vault user access ID format: %s (expected <VAULT ID>/<USER ID>)", id)
	}

	return s[0], s[1], nil
}
