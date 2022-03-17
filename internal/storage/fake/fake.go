package fake

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"sync"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage"
)

type repository struct {
	fakeFilePath         string
	usersByID            map[string]model.User
	groupsByID           map[string]model.Group
	membershipByID       map[string]model.Membership
	vaultsByID           map[string]model.Vault
	vaultGroupAccessByID map[string]model.VaultGroupAccess
	storageMu            sync.RWMutex
}

func NewRepository(fakeFilePath string) (storage.Repository, error) {
	// Try loading state from disk.
	// Ignore if file doesn't exists, it means its new storage.
	fks, _ := loadStorage(fakeFilePath)

	// Initialize storage.
	users := map[string]model.User{}
	if fks != nil && fks.Users != nil {
		users = fks.Users
	}

	groups := map[string]model.Group{}
	if fks != nil && fks.Groups != nil {
		groups = fks.Groups
	}

	members := map[string]model.Membership{}
	if fks != nil && fks.Groups != nil {
		members = fks.Members
	}

	vaults := map[string]model.Vault{}
	if fks != nil && fks.Groups != nil {
		vaults = fks.Vaults
	}

	vaultGroupAccess := map[string]model.VaultGroupAccess{}
	if fks != nil && fks.VaultGroupAccess != nil {
		vaultGroupAccess = fks.VaultGroupAccess
	}

	return &repository{
		fakeFilePath:         fakeFilePath,
		usersByID:            users,
		groupsByID:           groups,
		membershipByID:       members,
		vaultsByID:           vaults,
		vaultGroupAccessByID: vaultGroupAccess,
	}, nil
}

func (r *repository) CreateUser(ctx context.Context, user model.User) (*model.User, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := user.Email
	_, ok := r.usersByID[id]
	if ok {
		return nil, fmt.Errorf("user already exists")
	}

	user.ID = id
	r.usersByID[user.ID] = user

	err := r.dumpStorage()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) GetUserByID(ctx context.Context, id string) (*model.User, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	user, ok := r.usersByID[id]
	if !ok {
		return nil, fmt.Errorf("user does not exists")
	}

	return &user, nil
}

func (r *repository) GetUserByEmail(ctx context.Context, email string) (*model.User, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	// Fake storage doesn't need optimization.
	for _, u := range r.usersByID {
		if u.Email == email {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("user does not exists")
}

func (r *repository) EnsureUser(ctx context.Context, user model.User) (*model.User, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.usersByID[user.ID]
	if !ok {
		return nil, fmt.Errorf("user doesn't exists")
	}

	r.usersByID[user.Email] = user

	err := r.dumpStorage()
	if err != nil {
		return nil, err
	}

	return &user, nil
}

func (r *repository) DeleteUser(ctx context.Context, id string) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.usersByID[id]
	if !ok {
		return fmt.Errorf("user doesn't exists")
	}

	delete(r.usersByID, id)

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) CreateGroup(ctx context.Context, group model.Group) (*model.Group, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := group.Name
	_, ok := r.groupsByID[id]
	if ok {
		return nil, fmt.Errorf("group already exists")
	}

	group.ID = id
	r.groupsByID[group.ID] = group

	err := r.dumpStorage()
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *repository) GetGroupByID(ctx context.Context, id string) (*model.Group, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	group, ok := r.groupsByID[id]
	if !ok {
		return nil, fmt.Errorf("group does not exists")
	}

	return &group, nil
}

func (r *repository) GetGroupByName(ctx context.Context, name string) (*model.Group, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	// Fake storage doesn't need optimization.
	for _, u := range r.groupsByID {
		if u.Name == name {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("group does not exists")
}

func (r *repository) EnsureGroup(ctx context.Context, group model.Group) (*model.Group, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.groupsByID[group.ID]
	if !ok {
		return nil, fmt.Errorf("group doesn't exists")
	}

	r.groupsByID[group.Name] = group

	err := r.dumpStorage()
	if err != nil {
		return nil, err
	}

	return &group, nil
}

func (r *repository) DeleteGroup(ctx context.Context, id string) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.groupsByID[id]
	if !ok {
		return fmt.Errorf("group doesn't exists")
	}

	delete(r.groupsByID, id)

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) getMembershipID(groupID, userID string) string {
	return groupID + "/" + userID
}

func (r *repository) EnsureMembership(ctx context.Context, membership model.Membership) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := r.getMembershipID(membership.GroupID, membership.UserID)
	r.membershipByID[id] = membership

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteMembership(ctx context.Context, membership model.Membership) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := r.getMembershipID(membership.GroupID, membership.UserID)

	_, ok := r.membershipByID[id]
	if !ok {
		return fmt.Errorf("membership doesn't exists")
	}

	delete(r.membershipByID, id)

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetMembershipByID(ctx context.Context, groupID, userID string) (*model.Membership, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	id := r.getMembershipID(groupID, userID)
	m, ok := r.membershipByID[id]
	if !ok {
		return nil, fmt.Errorf("membership doesn't exists")
	}

	return &m, nil
}

func (r *repository) CreateVault(ctx context.Context, vault model.Vault) (*model.Vault, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := vault.Name
	_, ok := r.groupsByID[id]
	if ok {
		return nil, fmt.Errorf("vault already exists")
	}

	vault.ID = id
	r.vaultsByID[vault.ID] = vault

	err := r.dumpStorage()
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func (r *repository) GetVaultByID(ctx context.Context, id string) (*model.Vault, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	vault, ok := r.vaultsByID[id]
	if !ok {
		return nil, fmt.Errorf("vault does not exists")
	}

	return &vault, nil
}

func (r *repository) GetVaultByName(ctx context.Context, name string) (*model.Vault, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	// Fake storage doesn't need optimization.
	for _, u := range r.vaultsByID {
		if u.Name == name {
			return &u, nil
		}
	}

	return nil, fmt.Errorf("vault does not exists")
}

func (r *repository) EnsureVault(ctx context.Context, vault model.Vault) (*model.Vault, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.vaultsByID[vault.ID]
	if !ok {
		return nil, fmt.Errorf("vault doesn't exists")
	}

	r.vaultsByID[vault.Name] = vault

	err := r.dumpStorage()
	if err != nil {
		return nil, err
	}

	return &vault, nil
}

func (r *repository) DeleteVault(ctx context.Context, id string) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.vaultsByID[id]
	if !ok {
		return fmt.Errorf("vault doesn't exists")
	}

	delete(r.vaultsByID, id)

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) getVaultGroupAccessID(vaultID, groupID string) string {
	return vaultID + "/" + groupID
}

func (r *repository) EnsureVaultGroupAccess(ctx context.Context, groupAccess model.VaultGroupAccess) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := r.getVaultGroupAccessID(groupAccess.VaultID, groupAccess.GroupID)
	r.vaultGroupAccessByID[id] = groupAccess

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) DeleteVaultGroupAccess(ctx context.Context, vaultID string, groupID string) error {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	id := r.getVaultGroupAccessID(vaultID, groupID)

	_, ok := r.vaultGroupAccessByID[id]
	if !ok {
		return fmt.Errorf("vault access doesn't exists")
	}

	delete(r.vaultGroupAccessByID, id)

	err := r.dumpStorage()
	if err != nil {
		return err
	}

	return nil
}

func (r *repository) GetVaultGroupAccessByID(ctx context.Context, vaultID string, groupID string) (*model.VaultGroupAccess, error) {
	r.storageMu.RLock()
	defer r.storageMu.RUnlock()

	id := r.getVaultGroupAccessID(vaultID, groupID)
	v, ok := r.vaultGroupAccessByID[id]
	if !ok {
		return nil, fmt.Errorf("vault access doesn't exists")
	}

	return &v, nil
}

type fakeStorage struct {
	Users            map[string]model.User
	Groups           map[string]model.Group
	Members          map[string]model.Membership
	Vaults           map[string]model.Vault
	VaultGroupAccess map[string]model.VaultGroupAccess
}

func (r *repository) dumpStorage() error {
	fks := fakeStorage{
		Users:            r.usersByID,
		Groups:           r.groupsByID,
		Members:          r.membershipByID,
		Vaults:           r.vaultsByID,
		VaultGroupAccess: r.vaultGroupAccessByID,
	}

	data, err := json.MarshalIndent(fks, "", "\t")
	if err != nil {
		return fmt.Errorf("could not marshal storage: %w", err)
	}

	err = os.WriteFile(r.fakeFilePath, data, 0644)
	if err != nil {
		return fmt.Errorf("could not write file: %w", err)
	}

	return nil
}

func loadStorage(filePath string) (*fakeStorage, error) {
	fks := &fakeStorage{}

	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("could not read file: %w", err)
	}

	err = json.Unmarshal(data, fks)
	if err != nil {
		return nil, fmt.Errorf("could not unmarshal storage: %w", err)
	}

	return fks, nil
}
