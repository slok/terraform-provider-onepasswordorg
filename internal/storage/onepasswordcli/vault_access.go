package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r *Repository) EnsureVaultGroupAccess(ctx context.Context, groupAccess model.VaultGroupAccess) error {
	// Delete first to revoke every access and then add the required one again.
	// Ignore error in case there isn't an access already.
	_ = r.DeleteVaultGroupAccess(ctx, groupAccess.VaultID, groupAccess.GroupID)

	ps := mapModelToOpPermissions(groupAccess.Permissions)
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().GroupArg().GrantArg().VaultFlag(groupAccess.VaultID).GroupFlag(groupAccess.GroupID).NoInputFlag().PermissionsFlag(ps)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

func (r *Repository) DeleteVaultGroupAccess(ctx context.Context, vaultID string, groupID string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().GroupArg().RevokeArg().VaultFlag(vaultID).GroupFlag(groupID)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

func (r *Repository) GetVaultGroupAccessByID(ctx context.Context, vaultID string, groupID string) (*model.VaultGroupAccess, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().GroupArg().ListArg().RawStrArg(vaultID).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	// List vault groups and get the correct group.
	accesses := []opVaultGroupAccess{}
	err = json.Unmarshal([]byte(stdout), &accesses)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}
	var access *opVaultGroupAccess
	for _, a := range accesses {
		if a.GroupID == groupID {
			access = &a
			break
		}
	}

	if access == nil {
		return nil, fmt.Errorf("group access %q in vault %q not found", groupID, vaultID)
	}

	return &model.VaultGroupAccess{
		VaultID:     vaultID,
		GroupID:     groupID,
		Permissions: mapOpToModelPermissions(access.Permissions),
	}, nil
}

const (
	accessPermAllowViewing         = "allow_viewing"
	accessPermAllowEditing         = "allow_editing"
	accessPermAllowManaging        = "allow_managing"
	accessPermViewItems            = "view_items"
	accessPermCreateItems          = "create_items"
	accessPermEditItems            = "edit_items"
	accessPermArchiveItems         = "archive_items"
	accessPermDeleteItems          = "delete_items"
	accessPermViewAndCopyPasswords = "view_and_copy_passwords"
	accessPermViewItemHistory      = "view_item_history"
	accessPermImportItems          = "import_items"
	accessPermExportItems          = "export_items"
	accessPermCopyAndShareItems    = "copy_and_share_items"
	accessPermPrintItems           = "print_items"
	accessPermManageVault          = "manage_vault"
)

func mapModelToOpPermissions(p model.AccessPermissions) []string {
	ps := []string{}

	if p.AllowViewing {
		ps = append(ps, accessPermAllowViewing)
	}
	if p.AllowEditing {
		ps = append(ps, accessPermAllowEditing)
	}
	if p.AllowManaging {
		ps = append(ps, accessPermAllowManaging)
	}
	if p.ViewItems {
		ps = append(ps, accessPermViewItems)
	}
	if p.CreateItems {
		ps = append(ps, accessPermCreateItems)
	}
	if p.EditItems {
		ps = append(ps, accessPermEditItems)
	}
	if p.ArchiveItems {
		ps = append(ps, accessPermArchiveItems)
	}
	if p.DeleteItems {
		ps = append(ps, accessPermDeleteItems)
	}
	if p.ViewAndCopyPasswords {
		ps = append(ps, accessPermViewAndCopyPasswords)
	}
	if p.ViewItemHistory {
		ps = append(ps, accessPermViewItemHistory)
	}
	if p.ImportItems {
		ps = append(ps, accessPermImportItems)
	}
	if p.ExportItems {
		ps = append(ps, accessPermExportItems)
	}
	if p.CopyAndShareItems {
		ps = append(ps, accessPermCopyAndShareItems)
	}
	if p.PrintItems {
		ps = append(ps, accessPermPrintItems)
	}
	if p.ManageVault {
		ps = append(ps, accessPermManageVault)
	}

	return ps
}

func mapOpToModelPermissions(permissions []string) model.AccessPermissions {
	ap := model.AccessPermissions{}
	for _, p := range permissions {
		switch p {
		case accessPermAllowViewing:
			ap.AllowViewing = true
		case accessPermAllowEditing:
			ap.AllowEditing = true
		case accessPermAllowManaging:
			ap.AllowManaging = true
		case accessPermViewItems:
			ap.ViewItems = true
		case accessPermCreateItems:
			ap.CreateItems = true
		case accessPermEditItems:
			ap.EditItems = true
		case accessPermArchiveItems:
			ap.ArchiveItems = true
		case accessPermDeleteItems:
			ap.DeleteItems = true
		case accessPermViewAndCopyPasswords:
			ap.ViewAndCopyPasswords = true
		case accessPermViewItemHistory:
			ap.ViewItemHistory = true
		case accessPermImportItems:
			ap.ImportItems = true
		case accessPermExportItems:
			ap.ExportItems = true
		case accessPermCopyAndShareItems:
			ap.CopyAndShareItems = true
		case accessPermPrintItems:
			ap.PrintItems = true
		case accessPermManageVault:
			ap.ManageVault = true
		}
	}

	return ap
}

type opVaultGroupAccess struct {
	GroupID     string   `json:"id"`
	Permissions []string `json:"permissions"`
}
