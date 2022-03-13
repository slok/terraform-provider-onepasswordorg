package storage

import (
	"context"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	EnsureUser(ctx context.Context, user model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error
}
