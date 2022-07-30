package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

type resourceVaultUserAccessType struct{}

func (r resourceVaultUserAccessType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides vault access for a user.
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"vault_id": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The vault ID.",
			},
			"user_id": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The user ID.",
			},
			"permissions": permissionsAttribute,
		},
	}, nil
}

func (r resourceVaultUserAccessType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceVaultUserAccess{
		p: *prv,
	}, nil
}

type resourceVaultUserAccess struct {
	p provider
}

func (r resourceVaultUserAccess) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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

	err = r.p.repo.EnsureVaultUserAccess(ctx, *v)
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

func (r resourceVaultUserAccess) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVaultUserAccess VaultUserAccess
	diags := req.State.Get(ctx, &tfVaultUserAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get access.
	id := tfVaultUserAccess.ID.Value
	vaultID, userID, err := unpackVaultUserAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	access, err := r.p.repo.GetVaultUserAccessByID(ctx, vaultID, userID)
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

func (r resourceVaultUserAccess) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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

	err = r.p.repo.EnsureVaultUserAccess(ctx, *v)
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

func (r resourceVaultUserAccess) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVaultUserAccess VaultUserAccess
	diags := req.State.Get(ctx, &tfVaultUserAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete resource.
	id := tfVaultUserAccess.ID.Value
	vaultID, userID, err := unpackVaultUserAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	err = r.p.repo.DeleteVaultUserAccess(ctx, vaultID, userID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting access", fmt.Sprintf("Could not delete access %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceVaultUserAccess) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTfToModelVaultUserAccess(m VaultUserAccess) (*model.VaultUserAccess, error) {
	userID := m.UserID.Value
	vaultID := m.VaultID.Value

	// Check the ID is correct.
	if m.ID.Value != "" {
		vid, uid, err := unpackVaultUserAccessID(m.ID.Value)
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
		ID:          types.String{Value: id},
		UserID:      types.String{Value: m.UserID},
		VaultID:     types.String{Value: m.VaultID},
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
