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
				expCmd := `create user test@test.io  Test00`
				stdout := `{"uuid":"1234567890","createdAt":"2022-03-13T18:49:59Z","updatedAt":"2022-03-13T18:49:59Z","lastAuthAt":"0001-01-01T00:00:00Z","email":"test@test.io","firstName":"Test00","lastName":"","name":"Test00"}`
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
				expCmd := `create user test@test.io  Test00`
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
				expCmd := `get user test-id`
				stdout := `{"uuid":"1234567890","createdAt":"2021-09-08T07:45:22Z","updatedAt":"2021-09-08T07:47:02Z","lastAuthAt":"2022-03-12T14:23:17Z","email":"test@test.io","firstName":"Test00","lastName":"","name":"Test00","attrVersion":3,"keysetVersion":4,"state":"A","type":"R","avatar":"","language":"en","accountKeyFormat":"","accountKeyUuid":"","combinedPermissions":98765432}`
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
				expCmd := `get user test-id`
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
				expCmd := `get user test@test.io`
				stdout := `{"uuid":"1234567890","createdAt":"2021-09-08T07:45:22Z","updatedAt":"2021-09-08T07:47:02Z","lastAuthAt":"2022-03-12T14:23:17Z","email":"test@test.io","firstName":"Test00","lastName":"","name":"Test00","attrVersion":3,"keysetVersion":4,"state":"A","type":"R","avatar":"","language":"en","accountKeyFormat":"","accountKeyUuid":"","combinedPermissions":98765432}`
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
				expCmd := `get user test@test.io`
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
				expCmd := `edit user test-id --name Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
			expUser: &model.User{ID: "test-id", Email: "test@test.io", Name: "Test00"},
		},

		"Having an error while calling the op CLI, should fail.": {
			user: model.User{ID: "test-id", Email: "test@test.io", Name: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `edit user test-id --name Test00`
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
				expCmd := `delete user test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `delete user test-id`
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
