package provider

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-framework-validators/stringvalidator"
	"github.com/hashicorp/terraform-plugin-framework/path"
	"github.com/hashicorp/terraform-plugin-framework/resource"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/planmodifier"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"

	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

var (
	_ resource.Resource                = &userResource{}
	_ resource.ResourceWithConfigure   = &userResource{}
	_ resource.ResourceWithImportState = &userResource{}
)

func NewUserResource() resource.Resource {
	return &userResource{}
}

type userResource struct {
	repo storage.Repository
}

func (r *userResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_user"
}

func (r *userResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a User resource.

When a 1password user resources is created, it will be invited  by email.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
			},
			"name": schema.StringAttribute{
				Required: true,
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "The name of the user.",
			},
			"email": schema.StringAttribute{
				Required: true,
				PlanModifiers: []planmodifier.String{
					stringplanmodifier.RequiresReplace(),
				},
				Validators: []validator.String{
					stringvalidator.LengthAtLeast(1),
				},
				Description: "The email of the user.",
			},
		},
	}
}

func (r *userResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r *userResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var tfUser User
	diags := req.Plan.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create user.
	u := mapTfToModelUser(tfUser)
	newUser, err := r.repo.CreateUser(ctx, u)
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

func (r *userResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan.
	var tfUser User
	diags := req.State.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get user.
	id := tfUser.ID.ValueString()
	user, err := r.repo.GetUserByID(ctx, id)
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

func (r *userResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
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
	u.ID = state.ID.ValueString()

	newUser, err := r.repo.EnsureUser(ctx, u)
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

func (r *userResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan.
	var tfUser User
	diags := req.State.Get(ctx, &tfUser)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete user.
	id := tfUser.ID.ValueString()
	err := r.repo.DeleteUser(ctx, id)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting user", fmt.Sprintf("Could not delete user %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r *userResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}
