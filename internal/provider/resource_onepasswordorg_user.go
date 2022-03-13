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

type resourceUserType struct{}

func (r resourceUserType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides a User resource.

When a 1password user resources is created, it will be invited  by email.
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {

				Type:     types.StringType,
				Computed: true,
			},
			"name": {
				Type:        types.StringType,
				Required:    true,
				Description: "The name of the user.",
			},
			"email": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Description:   "The email of the user.",
			},
		},
	}, nil
}

func (r resourceUserType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceUser{
		p: *prv,
	}, nil
}

type resourceUser struct {
	p provider
}

func (r resourceUser) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfUser User
	diags := req.Plan.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create user.
	u := mapTfToModelUser(tfUser)
	newUser, err := r.p.repo.CreateUser(ctx, u)
	if err != nil {
		resp.Diagnostics.AddError("Error creating user", "Could not create user, unexpected error: "+err.Error())
		return
	}

	// Map user to tf model.
	newTfUser := mapModelToTfUser(*newUser)

	diags = resp.State.Set(ctx, newTfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceUser) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfUser User
	diags := req.State.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get user.
	id := tfUser.ID.Value
	user, err := r.p.repo.GetUserByID(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", fmt.Sprintf("Could not get user %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map user to tf model.
	readTfUser := mapModelToTfUser(*user)

	diags = resp.State.Set(ctx, readTfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceUser) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Get plan values.
	var plan User
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get current state.
	var state User
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan user as the new data and set ID from state.
	u := mapTfToModelUser(plan)
	u.ID = state.ID.Value

	newUser, err := r.p.repo.EnsureUser(ctx, u)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", "Could not update user, unexpected error: "+err.Error())
		return
	}

	// Map user to tf model.
	readTfUser := mapModelToTfUser(*newUser)

	diags = resp.State.Set(ctx, readTfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r resourceUser) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfUser User
	diags := req.State.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get user.
	id := tfUser.ID.Value
	err := r.p.repo.DeleteUser(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error reading user", fmt.Sprintf("Could not get user %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceUser) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

func mapTfToModelUser(u User) model.User {
	return model.User{
		ID:    u.ID.Value,
		Email: u.Email.Value,
		Name:  u.Name.Value,
	}
}

func mapModelToTfUser(u model.User) User {
	return User{
		ID:    types.String{Value: u.ID},
		Email: types.String{Value: u.Email},
		Name:  types.String{Value: u.Name},
	}
}
