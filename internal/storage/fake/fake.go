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

	return &repository{
		fakeFilePath: fakeFilePath,
		usersByID:    users,
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
	r.usersByID[user.Email] = user

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID})
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

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID})
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

	err := dumpStorage(r.fakeFilePath, fakeStorage{Users: r.usersByID})
	if err != nil {
		return err
	}

	return nil
}

type fakeStorage struct {
	Users map[string]model.User
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