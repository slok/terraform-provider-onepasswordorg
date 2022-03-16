package onepasswordcli

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

func (r Repository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.UserArg().ProvisionArg().EmailFlag(user.Email).NameFlag(user.Name).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opUser{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotUser := mapOpToModelUser(ou)

	return &gotUser, nil
}

func (r Repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.UserArg().GetArg().RawStrArg(id).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opUser{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotUser := mapOpToModelUser(ou)

	return &gotUser, nil
}

func (r Repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.UserArg().GetArg().RawStrArg(email).FormatJSONFlag()

	stdout, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	ou := opUser{}
	err = json.Unmarshal([]byte(stdout), &ou)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal op cli stdout: %w", err)
	}

	gotUser := mapOpToModelUser(ou)

	return &gotUser, nil
}

func (r Repository) EnsureUser(ctx context.Context, user model.User) (*model.User, error) {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.UserArg().EditArg().RawStrArg(user.ID).NameFlag(user.Name)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return nil, fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return &user, nil
}

func (r Repository) DeleteUser(ctx context.Context, id string) error {
	cmdArgs := &onePasswordCliCmd{}
	cmdArgs.UserArg().DeleteArg().RawStrArg(id)

	_, stderr, err := r.cli.RunOpCmd(ctx, cmdArgs.GetArgs())
	if err != nil {
		return fmt.Errorf("op cli command failed: %w: %s", err, stderr)
	}

	return nil
}

type opUser struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
}

func mapOpToModelUser(u opUser) model.User {
	return model.User{
		ID:    u.ID,
		Email: u.Email,
		Name:  u.Name,
	}
}
