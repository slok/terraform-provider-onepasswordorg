package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r *Repository) EnsureVaultUserAccess(ctx context.Context, userAccess model.VaultUserAccess) error {
	// Delete first to revoke every access and then add the required one again.
	// Ignore error in case there isn't an access already.
	_ = r.DeleteVaultUserAccess(ctx, userAccess.VaultID, userAccess.UserID)

	ps := mapModelToOpPermissions(userAccess.Permissions)
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().UserArg().GrantArg().VaultFlag(userAccess.VaultID).UserFlag(userAccess.UserID).NoInputFlag().PermissionsFlag(ps)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

func (r *Repository) DeleteVaultUserAccess(ctx context.Context, vaultID string, userID string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().UserArg().RevokeArg().VaultFlag(vaultID).UserFlag(userID)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

func (r *Repository) GetVaultUserAccessByID(ctx context.Context, vaultID string, userID string) (*model.VaultUserAccess, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().UserArg().ListArg().RawStrArg(vaultID).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	// List vault users and get the correct user.
	accesses := []opVaultUserAccess{}
	err = json.Unmarshal([]byte(stdout), &accesses)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}
	var access *opVaultUserAccess
	for _, a := range accesses {
		if a.UserID == userID {
			access = &a
			break
		}
	}

	if access == nil {
		return nil, fmt.Errorf("user access %q in vault %q not found", userID, vaultID)
	}

	return &model.VaultUserAccess{
		VaultID:     vaultID,
		UserID:      userID,
		Permissions: mapOpToModelPermissions(access.Permissions),
	}, nil
}

type opVaultUserAccess struct {
	UserID      string   `json:"id"`
	Permissions []string `json:"permissions"`
}
