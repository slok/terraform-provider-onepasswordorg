package storage

import (
	"context"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
)

type Repository interface {
	CreateUser(ctx context.Context, user model.User) (*model.User, error)
	GetUserByID(ctx context.Context, id string) (*model.User, error)
	GetUserByEmail(ctx context.Context, email string) (*model.User, error)
	EnsureUser(ctx context.Context, user model.User) (*model.User, error)
	DeleteUser(ctx context.Context, id string) error

	CreateGroup(ctx context.Context, group model.Group) (*model.Group, error)
	GetGroupByID(ctx context.Context, id string) (*model.Group, error)
	EnsureGroup(ctx context.Context, group model.Group) (*model.Group, error)
	DeleteGroup(ctx context.Context, id string) error

	EnsureMembership(ctx context.Context, membership model.Membership) error
	DeleteMembership(ctx context.Context, membership model.Membership) error
	GetMembershipByID(ctx context.Context, groupID, userID string) (*model.Membership, error)
}
