package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

type opGroupMember struct {
	ID   string `json:"uuid"`
	Role string `json:"role"`
}

func (r Repository) EnsureMembership(ctx context.Context, membership model.Membership) error {
	role, err := mapModelToOpRole(membership.Role)
	if err != nil {
		return fmt.Errorf("could not map role: %w", err)
	}

	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.AddArg().UserArg().RawStrArg(membership.UserID).RawStrArg(membership.GroupID).RoleFlag(role)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	// 1password doesn't know to add a member to a group with a specific role, so we would need to:
	// - Add user to group.
	// - Change role.
	// By default 1passwords add users as members, so in case the role is other than that make a second
	// call always. We assume this tradeoff of extra call in these cases to avoid the need to apply
	// tf twice.
	if membership.Role != model.MembershipRoleMember {
		_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
		if err != nil {
			return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
		}
	}

	return nil
}

func (r Repository) GetMembershipByID(ctx context.Context, groupID, userID string) (*model.Membership, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ListArg().UsersArg().GroupFlag(groupID)

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	// List group members and get the user.
	members := []opGroupMember{}
	err = json.Unmarshal([]byte(stdout), &members)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}
	var member *opGroupMember
	for _, m := range members {
		if m.ID == userID {
			member = &m
			break
		}
	}

	if member == nil {
		return nil, fmt.Errorf("member %q in group %q not found", userID, groupID)
	}

	role, err := mapOpToModelRole(member.Role)
	if err != nil {
		return nil, fmt.Errorf("invalid role: %w", err)
	}

	return &model.Membership{
		UserID:  userID,
		GroupID: groupID,
		Role:    role,
	}, nil
}

func (r Repository) DeleteMembership(ctx context.Context, membership model.Membership) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.RemoveArg().UserArg().RawStrArg(membership.UserID).RawStrArg(membership.GroupID)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

func mapModelToOpRole(m model.MembershipRole) (string, error) {
	switch m {
	case model.MembershipRoleMember:
		return "member", nil
	case model.MembershipRoleManager:
		return "manager", nil
	}

	return "", fmt.Errorf("invalid role")
}

func mapOpToModelRole(role string) (model.MembershipRole, error) {
	switch strings.ToLower(role) {
	case "member":
		return model.MembershipRoleMember, nil
	case "manager":
		return model.MembershipRoleManager, nil
	default:
		return model.MembershipRoleMember, fmt.Errorf("invalid role")
	}
}
