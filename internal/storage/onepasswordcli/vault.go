package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r Repository) CreateVault(ctx context.Context, vault model.Vault) (*model.Vault, error) {
	// 1password allows multiple vaults with the same name, we add this to make sure
	// this doesn't happen.
	_, err := r.GetVaultByName(ctx, vault.Name)
	if err == nil {
		return nil, fmt.Errorf("vault with name %q already exists", vault.Name)
	}

	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().CreateArg().RawStrArg(vault.Name).DescriptionFlag(vault.Description).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ov := opVault{}
	err = json.Unmarshal([]byte(stdout), &ov)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotVault := mapOpToModeVault(ov)

	return &gotVault, nil
}

func (r Repository) GetVaultByID(ctx context.Context, id string) (*model.Vault, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().GetArg().RawStrArg(id).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ov := opVault{}
	err = json.Unmarshal([]byte(stdout), &ov)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotVault := mapOpToModeVault(ov)

	return &gotVault, nil
}

func (r Repository) ListVaultsByUser(ctx context.Context, userID string) (*[]model.Vault, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().ListArg().UserFlag(userID).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ov := []opVault{}
	err = json.Unmarshal([]byte(stdout), &ov)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotVault := []model.Vault{}
	for _, a := range ov {
		gotVault = append(gotVault, mapOpToModeVault(a))
	}

	return &gotVault, nil
}

func (r Repository) GetVaultByName(ctx context.Context, name string) (*model.Vault, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().GetArg().RawStrArg(name).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ov := opVault{}
	err = json.Unmarshal([]byte(stdout), &ov)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotVault := mapOpToModeVault(ov)

	return &gotVault, nil
}

func (r Repository) EnsureVault(ctx context.Context, vault model.Vault) (*model.Vault, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().EditArg().RawStrArg(vault.ID).DescriptionFlag(vault.Description).NameFlag(vault.Name)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return &vault, nil
}

func (r Repository) DeleteVault(ctx context.Context, id string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.VaultArg().DeleteArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

type opVault struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func mapOpToModeVault(v opVault) model.Vault {
	return model.Vault{
		ID:          v.ID,
		Name:        v.Name,
		Description: v.Description,
	}
}
