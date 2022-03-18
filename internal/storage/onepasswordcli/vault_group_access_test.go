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

func TestRepositoryEnsureVaultGroupAccess(t *testing.T) {
	tests := map[string]struct {
		access model.VaultGroupAccess
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Creating a group access correctly, should return the data with the ID.": {
			access: model.VaultGroupAccess{
				VaultID: "vault-00",
				GroupID: "group-00",
				Permissions: model.AccessPermissions{
					AllowEditing:      true,
					ExportItems:       true,
					AllowViewing:      true,
					CopyAndShareItems: true,
				},
			},
			mock: func(m *onepasswordclimock.OpCli) {
				// First revokes everything.
				expCmd := `vault group revoke --vault vault-00 --group group-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)

				// Grant again.
				expCmd = `vault group grant --vault vault-00 --group group-00 --no-input --permissions allow_viewing,allow_editing,export_items,copy_and_share_items`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the create op CLI action, should fail.": {
			access: model.VaultGroupAccess{VaultID: "vault-00", GroupID: "group-00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault group revoke --vault vault-00 --group group-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)

				expCmd = `vault group grant --vault vault-00 --group group-00 --no-input`
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

			err = repo.EnsureVaultGroupAccess(context.TODO(), test.access)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetVaultGroupAccessByID(t *testing.T) {
	stdout := `
	[
  {
    "id": "group-id",
    "permissions": ["manage_vault"]
  },
  {
    "id": "group-id-2",
    "permissions": ["manage_vault"]
  },
  {
    "id": "group-id-3",
    "permissions": [
      "view_items",
      "create_items",
      "edit_items"
    ]
  }
]
`
	tests := map[string]struct {
		vaultID   string
		groupID   string
		mock      func(m *onepasswordclimock.OpCli)
		expAccess *model.VaultGroupAccess
		expErr    bool
	}{
		"Getting an access correctly, should return the acess data.": {
			vaultID: "vault-00",
			groupID: "group-id-3",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault group list vault-00 --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expAccess: &model.VaultGroupAccess{
				VaultID: "vault-00",
				GroupID: "group-id-3",
				Permissions: model.AccessPermissions{
					ViewItems:   true,
					CreateItems: true,
					EditItems:   true,
				},
			},
		},

		"Getting a missing access should fail.": {
			vaultID: "vault-00",
			groupID: "group-id-4",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault group list vault-00 --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"Having an error while calling the op CLI, should fail.": {
			vaultID: "vault-00",
			groupID: "group-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault group list vault-00 --format json`
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

			gotAccess, err := repo.GetVaultGroupAccessByID(context.TODO(), test.vaultID, test.groupID)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expAccess, gotAccess)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteVaultGroupAccess(t *testing.T) {
	tests := map[string]struct {
		access model.VaultGroupAccess
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Delete a access correctly, should delete the access.": {
			access: model.VaultGroupAccess{VaultID: "vault-00", GroupID: "group-00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault group revoke --vault vault-00 --group group-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			access: model.VaultGroupAccess{VaultID: "vault-00", GroupID: "group-00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault group revoke --vault vault-00 --group group-00`
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

			err = repo.DeleteVaultGroupAccess(context.TODO(), test.access.VaultID, test.access.GroupID)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
