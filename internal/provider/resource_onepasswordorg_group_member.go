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
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringdefault"
	"github.com/hashicorp/terraform-plugin-framework/resource/schema/stringplanmodifier"
	"github.com/hashicorp/terraform-plugin-framework/schema/validator"
	"github.com/hashicorp/terraform-plugin-framework/types"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

var (
	_ resource.Resource                = &groupMemberResource{}
	_ resource.ResourceWithConfigure   = &groupMemberResource{}
	_ resource.ResourceWithImportState = &groupMemberResource{}
)

func NewGroupMemberResource() resource.Resource {
	return &groupMemberResource{}
}

type groupMemberResource struct {
	repo storage.Repository
}

func (r *groupMemberResource) Metadata(_ context.Context, req resource.MetadataRequest, resp *resource.MetadataResponse) {
	resp.TypeName = req.ProviderTypeName + "_group_member"
}

func (r *groupMemberResource) Schema(_ context.Context, _ resource.SchemaRequest, resp *resource.SchemaResponse) {
	resp.Schema = schema.Schema{
		Description: `
Provides a user and group membership.

A 1password group membership will make a user part of a group with a role on that group.
`,
		Attributes: map[string]schema.Attribute{
			"id": schema.StringAttribute{
				Computed: true,
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
			"role": schema.StringAttribute{
				Optional: true,
				Computed: true,
				Default:  stringdefault.StaticString(tfMemberRoleMember),
				Validators: []validator.String{
					stringvalidator.OneOf(tfMemberRoleMember, tfMemberRoleManager),
				},
				Description: "The role of the user on the group (can be `member` or `manager`, by default member).",
			},
		},
	}
}

func (r *groupMemberResource) Configure(_ context.Context, req resource.ConfigureRequest, resp *resource.ConfigureResponse) {
	appServices := getAppServicesFromResourceRequest(&req)
	if appServices == nil {
		return
	}

	r.repo = appServices.Repository
}

func (r *groupMemberResource) Create(ctx context.Context, req resource.CreateRequest, resp *resource.CreateResponse) {
	// Retrieve values from plan.
	var tfMember Member
	diags := req.Plan.Get(ctx, &tfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Create membership.
	m, err := mapTfToModelMembership(tfMember)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping member", "Could not map membership:"+err.Error())
		return
	}

	err = r.repo.EnsureMembership(ctx, *m)
	if err != nil {
		resp.Diagnostics.AddError("Error creating membership", "Could not create membership, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfMember, err := mapModelToTfMembership(*m)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping member", "Could not map membership:"+err.Error())
		return
	}

	// Set on state.
	diags = resp.State.Set(ctx, newTfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r groupMemberResource) Read(ctx context.Context, req resource.ReadRequest, resp *resource.ReadResponse) {
	// Retrieve values from plan.
	var tfMember Member
	diags := req.State.Get(ctx, &tfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get member.
	id := tfMember.ID.ValueString()
	groupID, userID, err := unpackGroupMemberID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting member ID", "Could not get member ID:"+err.Error())
		return
	}

	member, err := r.repo.GetMembershipByID(ctx, groupID, userID)
	if err != nil {
		resp.Diagnostics.AddError("Error reading membership", fmt.Sprintf("Could not get membership %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Map user to tf model.
	readTfMember, err := mapModelToTfMembership(*member)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping member", "Could not map membership:"+err.Error())
		return
	}

	diags = resp.State.Set(ctx, readTfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r groupMemberResource) Update(ctx context.Context, req resource.UpdateRequest, resp *resource.UpdateResponse) {
	// Retrieve values from plan.
	var plan Member
	diags := req.Plan.Get(ctx, &plan)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	var state Member
	diags = req.State.Get(ctx, &state)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Use plan user as the new data and set ID from state.
	plan.ID = state.ID

	// Update membership.
	m, err := mapTfToModelMembership(plan)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping member", "Could not map membership:"+err.Error())
		return
	}

	err = r.repo.EnsureMembership(ctx, *m)
	if err != nil {
		resp.Diagnostics.AddError("Error updating user", "Could not create user, unexpected error: "+err.Error())
		return
	}

	// Map to tf model.
	newTfMember, err := mapModelToTfMembership(*m)
	if err != nil {
		resp.Diagnostics.AddError("Error mapping member", "Could not map membership:"+err.Error())
		return
	}

	// Set on state.
	diags = resp.State.Set(ctx, newTfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}
}

func (r groupMemberResource) Delete(ctx context.Context, req resource.DeleteRequest, resp *resource.DeleteResponse) {
	// Retrieve values from plan.
	var tfMember Member
	diags := req.State.Get(ctx, &tfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete member.
	id := tfMember.ID.ValueString()
	groupID, userID, err := unpackGroupMemberID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting member ID", "Could not get member ID:"+err.Error())
		return
	}

	m := model.Membership{GroupID: groupID, UserID: userID}
	err = r.repo.DeleteMembership(ctx, m)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting member", fmt.Sprintf("Could not delete member %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r *groupMemberResource) ImportState(ctx context.Context, req resource.ImportStateRequest, resp *resource.ImportStateResponse) {
	resource.ImportStatePassthroughID(ctx, path.Root("id"), req, resp)
}

const (
	tfMemberRoleMember  = "member"
	tfMemberRoleManager = "manager"
)

func mapTfToModelMembership(m Member) (*model.Membership, error) {
	groupID := m.GroupID.ValueString()
	userID := m.UserID.ValueString()

	// Check the ID is correct.
	if m.ID.ValueString() != "" {
		gid, uid, err := unpackGroupMemberID(m.ID.ValueString())
		if err != nil {
			return nil, err
		}

		if gid != groupID {
			return nil, fmt.Errorf("resource id is wrong based on group ID")
		}

		if uid != userID {
			return nil, fmt.Errorf("resource id is wrong based on user ID")
		}
	}

	var role model.MembershipRole
	switch m.Role.ValueString() {
	case tfMemberRoleMember:
		role = model.MembershipRoleMember
	case tfMemberRoleManager:
		role = model.MembershipRoleManager
	default:
		return nil, fmt.Errorf("the role %q is invalid", m.Role.ValueString())
	}

	return &model.Membership{
		UserID:  userID,
		GroupID: groupID,
		Role:    role,
	}, nil
}

func mapModelToTfMembership(m model.Membership) (*Member, error) {
	id := packGroupMemberID(m.GroupID, m.UserID)

	var role string
	switch m.Role {
	case model.MembershipRoleMember:
		role = tfMemberRoleMember
	case model.MembershipRoleManager:
		role = tfMemberRoleManager
	default:
		return nil, fmt.Errorf("the role %q is invalid", m.Role)
	}

	return &Member{
		ID:      types.StringValue(id),
		GroupID: types.StringValue(m.GroupID),
		UserID:  types.StringValue(m.UserID),
		Role:    types.StringValue(role),
	}, nil
}

func packGroupMemberID(groupID, userID string) string {
	return groupID + "/" + userID
}

func unpackGroupMemberID(id string) (groupID, userID string, err error) {
	s := strings.SplitN(id, "/", 2)
	if len(s) != 2 {
		return "", "", fmt.Errorf(
			"invalid group member ID format: %s (expected <GROUP ID>/<USER ID>)", id)
	}

	return s[0], s[1], nil
}
