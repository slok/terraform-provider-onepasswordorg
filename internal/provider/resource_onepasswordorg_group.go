package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework/diag"
	"github.com/hashicorp/terraform-plugin-framework/tfsdk"
	"github.com/hashicorp/terraform-plugin-framework/types"
	"github.com/hashicorp/terraform-plugin-go/tftypes"
	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

type resourceGroupType struct{}

func (r resourceGroupType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides a Group resource.

A 1password group is like a team that can contain people and can be used to give access to vaults as a
group of users.
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:          types.StringType,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Required:      true,
				Description:   "The name of the group.",
			},
			"description": {
				Type:        types.StringType,
				Optional:    true,
				Description: "The description of the group.",
			},
		},
	}, nil
}

func (r resourceGroupType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceGroup{
		p: *prv,
	}, nil
}

type resourceGroup struct {
	p provider
}

func (r resourceGroup) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfGroup Group
	diags := req.Plan.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create group.
	g := mapTfToModelGroup(tfGroup)
	newGroup, err := r.p.repo.CreateGroup(ctx, g)
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

func (r resourceGroup) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfGroup Group
	diags := req.State.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group.
	id := tfGroup.ID.Value
	group, err := r.p.repo.GetGroupByID(ctx, id)
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

func (r resourceGroup) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Get plan values.
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
	u := mapTfToModelGroup(plan)
	u.ID = state.ID.Value

	newGroup, err := r.p.repo.EnsureGroup(ctx, u)
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

func (r resourceGroup) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfGroup Group
	diags := req.State.Get(ctx, &tfGroup)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get group.
	id := tfGroup.ID.Value
	err := r.p.repo.DeleteGroup(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting group", fmt.Sprintf("Could not delete group %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceGroup) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func mapTfToModelGroup(g Group) model.Group {
	return model.Group{
		ID:          g.ID.Value,
		Name:        g.Name.Value,
		Description: g.Description.Value,
	}
}

func mapModelToTfGroup(g model.Group) Group {
	return Group{
		ID:          types.String{Value: g.ID},
		Name:        types.String{Value: g.Name},
		Description: types.String{Value: g.Description},
	}
}
