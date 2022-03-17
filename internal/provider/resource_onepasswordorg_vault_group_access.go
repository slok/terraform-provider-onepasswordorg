package provider

import (
	"context"
	"fmt"
	"strings"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

type resourceVaultGroupAccessType struct{}

func (r resourceVaultGroupAccessType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides vault access for a group.
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
			"group_id": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The group ID.",
			},
			"permissions": permissionsAttribute,
		},
	}, nil
}

func (r resourceVaultGroupAccessType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceVaultGroupAccess{
		p: *prv,
	}, nil
}

type resourceVaultGroupAccess struct {
	p provider
}

func (r resourceVaultGroupAccess) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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

	err = r.p.repo.EnsureVaultGroupAccess(ctx, *v)
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

func (r resourceVaultGroupAccess) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVaultGroupAccess VaultGroupAccess
	diags := req.State.Get(ctx, &tfVaultGroupAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get access.
	id := tfVaultGroupAccess.ID.Value
	vaultID, groupID, err := unpackVaultGroupAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	access, err := r.p.repo.GetVaultGroupAccessByID(ctx, vaultID, groupID)
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

func (r resourceVaultGroupAccess) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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

	err = r.p.repo.EnsureVaultGroupAccess(ctx, *v)
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

func (r resourceVaultGroupAccess) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVaultGroupAccess VaultGroupAccess
	diags := req.State.Get(ctx, &tfVaultGroupAccess)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete resource.
	id := tfVaultGroupAccess.ID.Value
	vaultID, groupID, err := unpackVaultGroupAccessID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting access ID", "Could not get access ID:"+err.Error())
		return
	}

	err = r.p.repo.DeleteVaultGroupAccess(ctx, vaultID, groupID)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting access", fmt.Sprintf("Could not delete access %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceVaultGroupAccess) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func mapTfToModelVaultGroupAccess(m VaultGroupAccess) (*model.VaultGroupAccess, error) {
	groupID := m.GroupID.Value
	vaultID := m.VaultID.Value

	// Check the ID is correct.
	if m.ID.Value != "" {
		vid, gid, err := unpackVaultGroupAccessID(m.ID.Value)
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
		ID:          types.String{Value: id},
		GroupID:     types.String{Value: m.GroupID},
		VaultID:     types.String{Value: m.VaultID},
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
