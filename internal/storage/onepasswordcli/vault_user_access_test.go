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

func TestRepositoryEnsureVaultUserAccess(t *testing.T) {
	tests := map[string]struct {
		access model.VaultUserAccess
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Creating a user access correctly, should return the data with the ID.": {
			access: model.VaultUserAccess{
				VaultID: "vault-00",
				UserID:  "user-00",
				Permissions: model.AccessPermissions{
					AllowEditing:      true,
					ExportItems:       true,
					AllowViewing:      true,
					CopyAndShareItems: true,
				},
			},
			mock: func(m *onepasswordclimock.OpCli) {
				// First revokes everything.
				expCmd := `vault user revoke --vault vault-00 --user user-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)

				// Grant again.
				expCmd = `vault user grant --vault vault-00 --user user-00 --no-input --permissions allow_viewing,allow_editing,export_items,copy_and_share_items`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the create op CLI action, should fail.": {
			access: model.VaultUserAccess{VaultID: "vault-00", UserID: "user-00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault user revoke --vault vault-00 --user user-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)

				expCmd = `vault user grant --vault vault-00 --user user-00 --no-input`
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

			err = repo.EnsureVaultUserAccess(context.TODO(), test.access)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetVaultUserAccessByID(t *testing.T) {
	stdout := `
	[
  {
    "id": "user-id",
    "permissions": ["manage_vault"]
  },
  {
    "id": "user-id-2",
    "permissions": ["manage_vault"]
  },
  {
    "id": "user-id-3",
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
		userID    string
		mock      func(m *onepasswordclimock.OpCli)
		expAccess *model.VaultUserAccess
		expErr    bool
	}{
		"Getting an access correctly, should return the acess data.": {
			vaultID: "vault-00",
			userID:  "user-id-3",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault user list vault-00 --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expAccess: &model.VaultUserAccess{
				VaultID: "vault-00",
				UserID:  "user-id-3",
				Permissions: model.AccessPermissions{
					ViewItems:   true,
					CreateItems: true,
					EditItems:   true,
				},
			},
		},

		"Getting a missing access should fail.": {
			vaultID: "vault-00",
			userID:  "user-id-4",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault user list vault-00 --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"Having an error while calling the op CLI, should fail.": {
			vaultID: "vault-00",
			userID:  "user-id",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault user list vault-00 --format json`
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

			gotAccess, err := repo.GetVaultUserAccessByID(context.TODO(), test.vaultID, test.userID)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expAccess, gotAccess)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteVaultUserAccess(t *testing.T) {
	tests := map[string]struct {
		access model.VaultUserAccess
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Delete a access correctly, should delete the access.": {
			access: model.VaultUserAccess{VaultID: "vault-00", UserID: "user-00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault user revoke --vault vault-00 --user user-00`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			access: model.VaultUserAccess{VaultID: "vault-00", UserID: "user-00"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `vault user revoke --vault vault-00 --user user-00`
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

			err = repo.DeleteVaultUserAccess(context.TODO(), test.access.VaultID, test.access.UserID)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
