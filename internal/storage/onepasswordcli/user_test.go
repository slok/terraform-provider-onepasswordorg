package onepasswordcli_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli"
	"github.com/slok/terraform-provider-onepasswordorg/internal/storage/onepasswordcli/onepasswordclimock"
)

func TestRepositoryCreateUser(t *testing.T) {
	tests := map[string]struct {
		user    model.User
		mock    func(m *onepasswordclimock.OpCli)
		expUser *model.User
		expErr  bool
	}{
		"Creating a user correctly, should return the data with the ID.": {
			user: model.User{Email: "test@test.io", Name: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user provision --email test@test.io --name Test00 --format json`
				stdout := `{"id":"1234567890","email":"test@test.io","name":"Test00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expUser: &model.User{
				ID:    "1234567890",
				Email: "test@test.io",
				Name:  "Test00",
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			user: model.User{Email: "test@test.io", Name: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user provision --email test@test.io --name Test00 --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotUser, err := repo.CreateUser(context.TODO(), test.user)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expUser, gotUser)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetUserByID(t *testing.T) {
	tests := map[string]struct {
		id      string
		mock    func(m *onepasswordclimock.OpCli)
		expUser *model.User
		expErr  bool
	}{
		"Getting a user correctly, should return the user data.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user get test-id --format json`
				stdout := `{"id":"1234567890","email":"test@test.io","name":"Test00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expUser: &model.User{
				ID:    "1234567890",
				Email: "test@test.io",
				Name:  "Test00",
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user get test-id --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotUser, err := repo.GetUserByID(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expUser, gotUser)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetUserByEmail(t *testing.T) {
	tests := map[string]struct {
		email   string
		mock    func(m *onepasswordclimock.OpCli)
		expUser *model.User
		expErr  bool
	}{
		"Getting a user correctly, should return the user data.": {
			email: "test@test.io",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user get test@test.io --format json`
				stdout := `{"id":"1234567890","email":"test@test.io","name":"Test00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expUser: &model.User{
				ID:    "1234567890",
				Email: "test@test.io",
				Name:  "Test00",
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			email: "test@test.io",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user get test@test.io --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotUser, err := repo.GetUserByEmail(context.TODO(), test.email)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expUser, gotUser)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryEnsureUser(t *testing.T) {
	tests := map[string]struct {
		user    model.User
		mock    func(m *onepasswordclimock.OpCli)
		expUser *model.User
		expErr  bool
	}{
		"Updating a user correctly, should update the user data.": {
			user: model.User{ID: "test-id", Email: "test@test.io", Name: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user edit test-id --name Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
			expUser: &model.User{ID: "test-id", Email: "test@test.io", Name: "Test00"},
		},

		"Having an error while calling the op CLI, should fail.": {
			user: model.User{ID: "test-id", Email: "test@test.io", Name: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user edit test-id --name Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			gotUser, err := repo.EnsureUser(context.TODO(), test.user)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expUser, gotUser)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteUser(t *testing.T) {
	tests := map[string]struct {
		id     string
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Delete a user correctly, should return the user data.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user delete test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user delete test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			require := require.New(t)
			assert := assert.New(t)

			mc := &onepasswordclimock.OpCli{}
			test.mock(mc)

			repo, err := onepasswordcli.NewRepository(mc)
			require.NoError(err)

			err = repo.DeleteUser(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
