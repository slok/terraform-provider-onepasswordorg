package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider/attributeutils"
)

type resourceVaultType struct{}

func (r resourceVaultType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides a vault resource.
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The name of the vault.",
			},
			"description": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultValue(types.String{Value: "Managed by Terraform"})},
				Description:   "The description of the vault.",
			},
		},
	}, nil
}

func (r resourceVaultType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceVault{
		p: *prv,
	}, nil
}

type resourceVault struct {
	p provider
}

func (r resourceVault) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVault Vault
	diags := req.Plan.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create vault.
	v := mapTfToModelVault(tfVault)
	newVault, err := r.p.repo.CreateVault(ctx, v)
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

func (r resourceVault) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVault Vault
	diags := req.State.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get resource.
	id := tfVault.ID.Value
	vault, err := r.p.repo.GetVaultByID(ctx, id)
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

func (r resourceVault) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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
	v.ID = state.ID.Value

	newVault, err := r.p.repo.EnsureVault(ctx, v)
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

func (r resourceVault) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfVault Vault
	diags := req.State.Get(ctx, &tfVault)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete resource.
	id := tfVault.ID.Value
	err := r.p.repo.DeleteVault(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting vault", fmt.Sprintf("Could not delete vault %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceVault) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

func mapTfToModelVault(v Vault) model.Vault {
	return model.Vault{
		ID:          v.ID.Value,
		Name:        v.Name.Value,
		Description: v.Description.Value,
	}
}

func mapModelToTfVault(u model.Vault) Vault {
	return Vault{
		ID:          types.String{Value: u.ID},
		Name:        types.String{Value: u.Name},
		Description: types.String{Value: u.Description},
	}
}
