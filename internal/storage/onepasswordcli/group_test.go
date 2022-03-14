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

func TestRepositoryCreateGroup(t *testing.T) {
	tests := map[string]struct {
		group    model.Group
		mock     func(m *onepasswordclimock.OpCli)
		expGroup *model.Group
		expErr   bool
	}{
		"Creating a group correctly, should return the data with the ID.": {
			group: model.Group{Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get group test-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("group doesn't exist"))

				expCmd = `create group test-00  --description Test00`
				stdout := `{"uuid":"1234567890","type":"U","name":"test-00","desc":"Test00","createdAt":"2022-03-14T07:48:26.179385832+01:00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expGroup: &model.Group{
				ID:          "1234567890",
				Name:        "test-00",
				Description: "Test00",
			},
		},

		"Creating a group that already exists, should  fail.": {
			group: model.Group{Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get group test-00`
				stdout := `{"uuid":"1234567890","type":"U","name":"test-00","desc":"Test00","createdAt":"2022-03-14T07:48:26.179385832+01:00"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"Having an error while calling the create op CLI action, should fail.": {
			group: model.Group{Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get group test-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("group doesn't exist"))

				expCmd = `create group test-00  --description Test00`
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

			gotGroup, err := repo.CreateGroup(context.TODO(), test.group)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expGroup, gotGroup)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetGroupByID(t *testing.T) {
	tests := map[string]struct {
		id       string
		mock     func(m *onepasswordclimock.OpCli)
		expGroup *model.Group
		expErr   bool
	}{
		"Getting a group correctly, should return the group data.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get group test-id`
				stdout := `{"uuid":"1234567890","type":"U","name":"test-00","desc":"Test00","createdAt":"2022-03-14T06:48:26Z","updatedAt":"2022-03-14T06:48:26Z","activeKeysetUuid":"231321321","attrVersion":1,"state":"A","permissions":16,"pubKey":{"alg":"RSA-OAEP","kid":"3213213213","ext":true,"e":"AQAB","n":"432443","key_ops":["encrypt"],"kty":"RSA"}}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expGroup: &model.Group{
				ID:          "1234567890",
				Name:        "test-00",
				Description: "Test00",
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `get group test-id`
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

			gotGroup, err := repo.GetGroupByID(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expGroup, gotGroup)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryEnsureGroup(t *testing.T) {
	tests := map[string]struct {
		group    model.Group
		mock     func(m *onepasswordclimock.OpCli)
		expGroup *model.Group
		expErr   bool
	}{
		"Updating a group correctly, should update the user data.": {
			group: model.Group{ID: "test-id", Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `edit group test-id --description Test00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
			expGroup: &model.Group{ID: "test-id", Name: "test-00", Description: "Test00"},
		},

		"Having an error while calling the op CLI, should fail.": {
			group: model.Group{ID: "test-id", Name: "test-00", Description: "Test00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `edit group test-id --description Test00`
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

			gotUser, err := repo.EnsureGroup(context.TODO(), test.group)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expGroup, gotUser)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteGroup(t *testing.T) {
	tests := map[string]struct {
		id     string
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Delete a group correctly, should return the group data.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `delete group test-id`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "test-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `delete group test-id`
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

			err = repo.DeleteGroup(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
