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
	GetGroupByName(ctx context.Context, name string) (*model.Group, error)
	EnsureGroup(ctx context.Context, group model.Group) (*model.Group, error)
	DeleteGroup(ctx context.Context, id string) error

	CreateVault(ctx context.Context, vault model.Vault) (*model.Vault, error)
	GetVaultByID(ctx context.Context, id string) (*model.Vault, error)
	GetVaultByName(ctx context.Context, name string) (*model.Vault, error)
	EnsureVault(ctx context.Context, vault model.Vault) (*model.Vault, error)
	DeleteVault(ctx context.Context, id string) error

	EnsureMembership(ctx context.Context, membership model.Membership) error
	DeleteMembership(ctx context.Context, membership model.Membership) error
	GetMembershipByID(ctx context.Context, groupID, userID string) (*model.Membership, error)

	EnsureVaultGroupAccess(ctx context.Context, groupAccess model.VaultGroupAccess) error
	DeleteVaultGroupAccess(ctx context.Context, vaultID string, groupID string) error
	GetVaultGroupAccessByID(ctx context.Context, vaultID string, groupID string) (*model.VaultGroupAccess, error)

	EnsureVaultUserAccess(ctx context.Context, userAccess model.VaultUserAccess) error
	DeleteVaultUserAccess(ctx context.Context, vaultID string, userID string) error
	GetVaultUserAccessByID(ctx context.Context, vaultID string, userID string) (*model.VaultUserAccess, error)

	CreateItem(ctx context.Context, item model.Item) (*model.Item, error)
	GetItemByID(ctx context.Context, id string) (*model.Item, error)
	GetItemByTitle(ctx context.Context, vaultID string, title string) (*model.Item, error)
	EnsureItem(ctx context.Context, item model.Item) (*model.Item, error)
	DeleteItem(ctx context.Context, id string) error
}
