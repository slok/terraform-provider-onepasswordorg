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

type resourceGroupMemberType struct{}

func (r resourceGroupMemberType) GetSchema(_ context.Context) (tfsdk.Schema, diag.Diagnostics) {
	return tfsdk.Schema{
		Description: `
Provides a user and group membership.

A 1password group membership will make a user part of a group with a role on that group.
`,
		Attributes: map[string]tfsdk.Attribute{
			"id": {
				Type:     types.StringType,
				Computed: true,
			},
			"user_id": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The user ID.",
			},
			"group_id": {
				Type:          types.StringType,
				Required:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{tfsdk.RequiresReplace()},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The group ID.",
			},
			"role": {
				Type:          types.StringType,
				Optional:      true,
				Computed:      true,
				PlanModifiers: tfsdk.AttributePlanModifiers{attributeutils.DefaultString("member")},
				Validators:    []tfsdk.AttributeValidator{attributeutils.NonEmptyString},
				Description:   "The role of the user on the group (can be `member` or `manager`, by default member).",
			},
		},
	}, nil
}

func (r resourceGroupMemberType) NewResource(_ context.Context, p tfsdk.Provider) (tfsdk.Resource, diag.Diagnostics) {
	prv := p.(*provider)
	return resourceGroupMember{
		p: *prv,
	}, nil
}

type resourceGroupMember struct {
	p provider
}

func (r resourceGroupMember) Create(ctx context.Context, req tfsdk.CreateResourceRequest, resp *tfsdk.CreateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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

	err = r.p.repo.EnsureMembership(ctx, *m)
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

func (r resourceGroupMember) Read(ctx context.Context, req tfsdk.ReadResourceRequest, resp *tfsdk.ReadResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfMember Member
	diags := req.State.Get(ctx, &tfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Get member.
	id := tfMember.ID.Value
	groupID, userID, err := unpackGroupMemberID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting member ID", "Could not get member ID:"+err.Error())
		return
	}

	member, err := r.p.repo.GetMembershipByID(ctx, groupID, userID)
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

func (r resourceGroupMember) Update(ctx context.Context, req tfsdk.UpdateResourceRequest, resp *tfsdk.UpdateResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

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

	err = r.p.repo.EnsureMembership(ctx, *m)
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

func (r resourceGroupMember) Delete(ctx context.Context, req tfsdk.DeleteResourceRequest, resp *tfsdk.DeleteResourceResponse) {
	if !r.p.configured {
		resp.Diagnostics.AddError("Provider not configured", "The provider hasn't been configured before apply.")
		return
	}

	// Retrieve values from plan.
	var tfMember Member
	diags := req.State.Get(ctx, &tfMember)
	resp.Diagnostics.Append(diags...)
	if resp.Diagnostics.HasError() {
		return
	}

	// Delete member.
	id := tfMember.ID.Value
	groupID, userID, err := unpackGroupMemberID(id)
	if err != nil {
		resp.Diagnostics.AddError("Error getting member ID", "Could not get member ID:"+err.Error())
		return
	}

	m := model.Membership{GroupID: groupID, UserID: userID}
	err = r.p.repo.DeleteMembership(ctx, m)
	if err != nil {
		resp.Diagnostics.AddError("Error deleting member", fmt.Sprintf("Could not delete member %q, unexpected error: %s", id, err.Error()))
		return
	}

	// Remove resource from state.
	resp.State.RemoveResource(ctx)
}

func (r resourceGroupMember) ImportState(ctx context.Context, req tfsdk.ImportResourceStateRequest, resp *tfsdk.ImportResourceStateResponse) {
	// Save the import identifier in the id attribute.
	tfsdk.ResourceImportStatePassthroughID(ctx, tftypes.NewAttributePath().WithAttributeName("id"), req, resp)
}

const (
	tfMemberRoleMember  = "member"
	tfMemberRoleManager = "manager"
)

func mapTfToModelMembership(m Member) (*model.Membership, error) {
	groupID := m.GroupID.Value
	userID := m.UserID.Value

	// Check the ID is correct.
	if m.ID.Value != "" {
		gid, uid, err := unpackGroupMemberID(m.ID.Value)
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
	switch m.Role.Value {
	case tfMemberRoleMember:
		role = model.MembershipRoleMember
	case tfMemberRoleManager:
		role = model.MembershipRoleManager
	default:
		return nil, fmt.Errorf("the role %q is invalid", m.Role.Value)
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
		ID:      types.String{Value: id},
		GroupID: types.String{Value: m.GroupID},
		UserID:  types.String{Value: m.UserID},
		Role:    types.String{Value: role},
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
