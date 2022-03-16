package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r Repository) CreateGroup(ctx context.Context, group model.Group) (*model.Group, error) {
	// 1password allows multiple groups with the same name, we add this to make sure
	// this doesn't happen.
	_, err := r.GetGroupByName(ctx, group.Name)
	if err == nil {
		return nil, fmt.Errorf("group with name %q already exists", group.Name)
	}

	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.GroupArg().CreateArg().RawStrArg(group.Name).DescriptionFlag(group.Description).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	og := opGroup{}
	err = json.Unmarshal([]byte(stdout), &og)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotGroup := mapOpToModelGroup(og)

	return &gotGroup, nil
}
func (r Repository) GetGroupByID(ctx context.Context, id string) (*model.Group, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.GroupArg().GetArg().RawStrArg(id).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	og := opGroup{}
	err = json.Unmarshal([]byte(stdout), &og)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotGroup := mapOpToModelGroup(og)

	return &gotGroup, nil
}

func (r Repository) GetGroupByName(ctx context.Context, name string) (*model.Group, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.GroupArg().GetArg().RawStrArg(name).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	og := opGroup{}
	err = json.Unmarshal([]byte(stdout), &og)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotGroup := mapOpToModelGroup(og)

	return &gotGroup, nil
}

func (r Repository) EnsureGroup(ctx context.Context, group model.Group) (*model.Group, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.GroupArg().EditArg().RawStrArg(group.ID).DescriptionFlag(group.Description)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return &group, nil
}
func (r Repository) DeleteGroup(ctx context.Context, id string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.GroupArg().DeleteArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

type opGroup struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func mapOpToModelGroup(u opGroup) model.Group {
	return model.Group{
		ID:          u.ID,
		Name:        u.Name,
		Description: u.Description,
	}
}
