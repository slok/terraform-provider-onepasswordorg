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
	fakeFilePath string
	usersByID    map[string]model.User
	groupsByID   map[string]model.Group
	storageMu    sync.RWMutex
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

	return &repository{
		fakeFilePath: fakeFilePath,
		usersByID:    users,
		groupsByID:   groups,
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

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID, Groups: r.groupsByID})
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

func (r *repository) EnsureUser(ctx context.Context, user model.User) (*model.User, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.usersByID[user.ID]
	if !ok {
		return nil, fmt.Errorf("user doesn't exists")
	}

	r.usersByID[user.Email] = user

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID, Groups: r.groupsByID})
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

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID, Groups: r.groupsByID})
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

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID, Groups: r.groupsByID})
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
func (r *repository) EnsureGroup(ctx context.Context, group model.Group) (*model.Group, error) {
	r.storageMu.Lock()
	defer r.storageMu.Unlock()

	_, ok := r.groupsByID[group.ID]
	if !ok {
		return nil, fmt.Errorf("group doesn't exists")
	}

	r.groupsByID[group.Name] = group

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID, Groups: r.groupsByID})
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

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID, Groups: r.groupsByID})
	if err != nil {
		return err
	}

	return nil
}

type fakeStorage struct {
	Users  map[string]model.User
	Groups map[string]model.Group
}

func dumpStorage(filePath string, fks fakeStorage) error {
	data, err := json.Marshal(fks)
	if err != nil {
		return fmt.Errorf("could not marshal storage: %w", err)
	}

	err = os.WriteFile(filePath, data, 0644)
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
