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

func TestRepositoryCreateServiceAccount(t *testing.T) {
	tests := map[string]struct {
		sa     model.ServiceAccount
		mock   func(m *onepasswordclimock.OpCli)
		expSA  *model.ServiceAccount
		expErr bool
	}{
		"Creating a service account correctly should return the data including the token.": {
			sa: model.ServiceAccount{Name: "ci-bot"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account create ci-bot --can-create-vaults --format json`
				stdout := `{"id":"SAXXXX","name":"ci-bot","token":"ops_secret_token"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expSA: &model.ServiceAccount{
				ID:    "SAXXXX",
				Name:  "ci-bot",
				Token: "ops_secret_token",
			},
		},

		"When create response omits name (op CLI v2.34+), should fall back to input name.": {
			sa: model.ServiceAccount{Name: "ci-bot"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account create ci-bot --can-create-vaults --format json`
				stdout := `{"id":"SAXXXX","token":"ops_secret_token"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expSA: &model.ServiceAccount{
				ID:    "SAXXXX",
				Name:  "ci-bot",
				Token: "ops_secret_token",
			},
		},

		"When create response uses uuid instead of id (op CLI v2.34+), should resolve correctly.": {
			sa: model.ServiceAccount{Name: "ci-bot"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account create ci-bot --can-create-vaults --format json`
				stdout := `{"uuid":"SAXXXX","name":"ci-bot","token":"ops_secret_token"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expSA: &model.ServiceAccount{
				ID:    "SAXXXX",
				Name:  "ci-bot",
				Token: "ops_secret_token",
			},
		},

		"When create response omits both id and uuid, should return an error.": {
			sa: model.ServiceAccount{Name: "ci-bot"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account create ci-bot --can-create-vaults --format json`
				stdout := `{"name":"ci-bot","token":"ops_secret_token"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"When create response omits token, should return an error.": {
			sa: model.ServiceAccount{Name: "ci-bot"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account create ci-bot --can-create-vaults --format json`
				stdout := `{"id":"SAXXXX","name":"ci-bot"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expErr: true,
		},

		"Having an error while calling the op CLI, should fail.": {
			sa: model.ServiceAccount{Name: "ci-bot"},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account create ci-bot --can-create-vaults --format json`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", fmt.Errorf("something"))
			},
			expErr: true,
		},

		"Creating with vault_access should emit one --vault flag per vault, in sorted key order.": {
			sa: model.ServiceAccount{
				Name: "ci-bot",
				VaultAccess: map[string][]model.ServiceAccountVaultPermission{
					"zebra-vault": {model.ServiceAccountVaultPermissionReadItems},
					"alpha-vault": {
						model.ServiceAccountVaultPermissionReadItems,
						model.ServiceAccountVaultPermissionWriteItems,
					},
				},
			},
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := []string{
					"service-account", "create", "ci-bot",
					"--can-create-vaults", "--format", "json",
					"--vault", "alpha-vault:read_items,write_items",
					"--vault", "zebra-vault:read_items",
				}
				stdout := `{"id":"SAXXXX","name":"ci-bot","token":"ops_secret_token"}`
				m.On("RunOpCmd", mock.Anything, expCmd).Once().Return(stdout, "", nil)
			},
			expSA: &model.ServiceAccount{
				ID:    "SAXXXX",
				Name:  "ci-bot",
				Token: "ops_secret_token",
				VaultAccess: map[string][]model.ServiceAccountVaultPermission{
					"zebra-vault": {model.ServiceAccountVaultPermissionReadItems},
					"alpha-vault": {
						model.ServiceAccountVaultPermissionReadItems,
						model.ServiceAccountVaultPermissionWriteItems,
					},
				},
			},
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

			gotSA, err := repo.CreateServiceAccount(context.TODO(), test.sa)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expSA, gotSA)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryGetServiceAccountByID(t *testing.T) {
	tests := map[string]struct {
		id     string
		mock   func(m *onepasswordclimock.OpCli)
		expSA  *model.ServiceAccount
		expErr bool
	}{
		"Getting a service account by ID should confirm existence via ratelimit (without returning name).": {
			id: "SAXXXX",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account ratelimit SAXXXX`
				stdout := `{"uuid":"SAXXXX"}`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return(stdout, "", nil)
			},
			expSA: &model.ServiceAccount{
				ID: "SAXXXX",
			},
		},

		"Having an error while calling the op CLI (e.g. SA not found), should fail.": {
			id: "SAXXXX",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `service-account ratelimit SAXXXX`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "\"SAXXXX\" isn't a service account", fmt.Errorf("exit status 1"))
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

			gotSA, err := repo.GetServiceAccountByID(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else if assert.NoError(err) {
				assert.Equal(test.expSA, gotSA)
			}

			mc.AssertExpectations(t)
		})
	}
}

func TestRepositoryDeleteServiceAccount(t *testing.T) {
	tests := map[string]struct {
		id     string
		mock   func(m *onepasswordclimock.OpCli)
		expErr bool
	}{
		"Deleting a service account correctly should succeed.": {
			id: "SAXXXX",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user delete SAXXXX`
				m.On("RunOpCmd", mock.Anything, strings.Fields(expCmd)).Once().Return("", "", nil)
			},
		},

		"Having an error while calling the op CLI, should fail.": {
			id: "SAXXXX",
			mock: func(m *onepasswordclimock.OpCli) {
				expCmd := `user delete SAXXXX`
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

			err = repo.DeleteServiceAccount(context.TODO(), test.id)

			if test.expErr {
				assert.Error(err)
			} else {
				assert.NoError(err)
			}

			mc.AssertExpectations(t)
		})
	}
}
