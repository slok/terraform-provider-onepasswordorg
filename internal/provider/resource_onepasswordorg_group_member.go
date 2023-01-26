package provider

import (
	"strings"

	"context"
	"fmt"

	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/validation"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func resourceGroupMember() *schema.Resource {
	return &schema.Resource{
		Description: `
Provides a user and group membership.

A 1password group membership will make a user part of a group with a role on that group.
    `,
		CreateContext: resourceGroupMemberCreate,
		ReadContext:   resourceGroupMemberRead,
		UpdateContext: resourceGroupMemberUpdate,
		DeleteContext: resourceGroupMemberDelete,

		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"user_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The user ID.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"group_id": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				Description:  "The group ID.",
				ValidateFunc: validation.StringIsNotEmpty,
			},
			"role": {
				Type:         schema.TypeString,
				Optional:     true,
				Default:      "member",
				Description:  "The role of the user on the group (can be `member` or `manager`, by default member).",
				ValidateFunc: validation.StringIsNotEmpty,
			},
		},
	}
}

func resourceGroupMemberCreate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	m, err := dataToModelMembership(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = p.repo.EnsureMembership(ctx, *m)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	mapModelToDataMembership(*m, data)

	return diags
}

func resourceGroupMemberRead(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	id := data.Id()
	groupID, userID, err := unpackGroupMemberID(id)
	if err != nil {
		return diag.Errorf("Error getting member ID: " + "Could not get member ID:" + err.Error())
	}

	member, err := p.repo.GetMembershipByID(ctx, groupID, userID)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get group %q, unexpected error: %s", id, err.Error()))
	}

	mapModelToDataMembership(*member, data)
	return diags
}

func resourceGroupMemberUpdate(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	m, err := dataToModelMembership(data)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	err = p.repo.EnsureMembership(ctx, *m)
	if err != nil {
		return diag.Errorf(err.Error())
	}

	mapModelToDataMembership(*m, data)

	return diags
}

func resourceGroupMemberDelete(ctx context.Context, data *schema.ResourceData, meta interface{}) diag.Diagnostics {
	p := meta.(ProviderConfig)
	var diags diag.Diagnostics
	if !p.configured {
		return diag.Errorf("Provider not configured:" + "The provider hasn't been configured before apply.")
	}

	// Get group.
	// Delete member.
	id := data.Id()
	groupID, userID, err := unpackGroupMemberID(id)
	if err != nil {
		return diag.Errorf("Error getting member ID: " + "Could not get member ID:" + err.Error())
	}

	m := model.Membership{GroupID: groupID, UserID: userID}
	err = p.repo.DeleteMembership(ctx, m)
	if err != nil {
		return diag.Errorf("Error reading group:" + fmt.Sprintf("Could not get group %q, unexpected error: %s", id, err.Error()))
	}

	return diags
}

const (
	tfMemberRoleMember  = "member"
	tfMemberRoleManager = "manager"
)

func dataToModelMembership(data *schema.ResourceData) (*model.Membership, error) {
	groupID := data.Get("group_id").(string)
	userID := data.Get("user_id").(string)
	id := data.Id()

	// Check the ID is correct.
	if id != "" {
		gid, uid, err := unpackGroupMemberID(id)
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
	switch data.Get("role").(string) {
	case tfMemberRoleMember:
		role = model.MembershipRoleMember
	case tfMemberRoleManager:
		role = model.MembershipRoleManager
	default:
		return nil, fmt.Errorf("the role %q is invalid", data.Get("role").(string))
	}

	return &model.Membership{
		UserID:  userID,
		GroupID: groupID,
		Role:    role,
	}, nil
}

func mapModelToDataMembership(m model.Membership, data *schema.ResourceData) error {
	id := packGroupMemberID(m.GroupID, m.UserID)

	var role string
	switch m.Role {
	case model.MembershipRoleMember:
		role = tfMemberRoleMember
	case model.MembershipRoleManager:
		role = tfMemberRoleManager
	default:
		return fmt.Errorf("the role %q is invalid", m.Role)
	}

	data.SetId(id)
	data.Set("group_id", m.GroupID)
	data.Set("user_id", m.UserID)
	data.Set("role", role)
	return nil
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
