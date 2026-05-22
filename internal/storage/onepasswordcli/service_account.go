package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

// opServiceAccount is the JSON representation returned by the op CLI for service accounts.
// The Token field is only populated on creation.
// op CLI v2.34+ uses "uuid" instead of "id"; we accept both.
type opServiceAccount struct {
	ID    string `json:"id"`
	UUID  string `json:"uuid"`
	Name  string `json:"name"`
	Token string `json:"token,omitempty"`
}

// resolveID returns the service account identifier, preferring "id" but
// falling back to "uuid" for op CLI v2.34+.
func (o opServiceAccount) resolveID() string {
	if o.ID != "" {
		return o.ID
	}
	return o.UUID
}

func (r Repository) CreateServiceAccount(ctx context.Context, sa model.ServiceAccount) (*model.ServiceAccount, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ServiceAccountArg().CreateArg().RawStrArg(sa.Name).CanCreateVaultsFlag().FormatJSONFlag()

	// Vault access is configurable only at creation time via repeated `--vault`
	// flags. Iterate in sorted-key order so command construction is
	// deterministic (helps tests and debug output).
	vaultNames := make([]string, 0, len(sa.VaultAccess))
	for name := range sa.VaultAccess {
		vaultNames = append(vaultNames, name)
	}
	sort.Strings(vaultNames)
	for _, name := range vaultNames {
		cmdArgs.ServiceAccountVaultFlag(name, sa.VaultAccess[name])
	}

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	var created opServiceAccount
	if err := json.Unmarshal([]byte(stdout), &created); err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w: raw output: %s", err, stdout)
	}

	id := created.resolveID()
	if id == "" {
		return nil, fmt.Errorf("service account create response missing id/uuid, raw output: %s", stdout)
	}
	if created.Token == "" {
		return nil, fmt.Errorf("service account create response missing token, raw output: %s", stdout)
	}
	// op CLI v2.34+ omits name from the create response; fall back to the input name.
	if created.Name == "" {
		created.Name = sa.Name
	}

	return &model.ServiceAccount{
		ID:          id,
		Name:        created.Name,
		Token:       created.Token,
		VaultAccess: sa.VaultAccess,
	}, nil
}

func (r Repository) GetServiceAccountByID(ctx context.Context, id string) (*model.ServiceAccount, error) {
	// op CLI v2.34+ removed the 'service-account get' subcommand. Use ratelimit as an
	// existence probe: it exits non-zero if the service account does not exist.
	// The ratelimit response does not include the service account name; the resource's
	// Read function preserves it from prior Terraform state.
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.ServiceAccountArg().RatelimitArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return &model.ServiceAccount{ID: id}, nil
}

func (r Repository) GetServiceAccountByName(ctx context.Context, name string) (*model.ServiceAccount, error) {
	return r.GetServiceAccountByID(ctx, name)
}

func (r Repository) DeleteServiceAccount(ctx context.Context, id string) error {
	// op CLI v2.34's `service-account` subcommand only exposes `create` and
	// `ratelimit`; there is no `service-account delete`. Service accounts
	// appear as users in `op user list`, so deletion goes through
	// `op user delete <SA-ID>`.
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.UserArg().DeleteArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}
