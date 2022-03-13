package onepasswordcli

import (
	"context"
	"fmt"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

// OpCli knows how to execute Op CLI commands.
type OpCli interface {
	RunOpCmd(ctx context.Context, args []string) (stdout, stderr string, err error)
}

type opCli struct {
	sessionToken string
}

func (o opCli) RunOpCmd(ctx context.Context, args []string) (stdout, stderr string, err error) {
	if o.sessionToken == "" {
		return "", "", fmt.Errorf("unauthenticated, op cli must singin first")
	}

	return "", "", fmt.Errorf("not implemented")
}

// NewOpCLI creates a new signed OpCLI command executor.
func NewOpCli(address, email, secretKey string) (OpCli, error) {
	return opCli{sessionToken: ""}, nil
}

// NewRepository returns a 1password CLI (op) based respoitory.
func NewRepository(cli OpCli) (*Repository, error) {
	return &Repository{}, nil
}

// Repository knows how to execute 1password operations using 1password CLI.
type Repository struct{}

func (r Repository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	return nil, nil
}

func (r Repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	return nil, nil
}

func (r Repository) EnsureUser(ctx context.Context, user model.User) (*model.User, error) {
	return nil, nil
}

func (r Repository) DeleteUser(ctx context.Context, id string) error {
	return nil
}
