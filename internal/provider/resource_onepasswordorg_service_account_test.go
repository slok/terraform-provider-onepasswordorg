package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccServiceAccountCreateDelete checks that a service account is created and deleted.
func TestAccServiceAccountCreateDelete(t *testing.T) {
	tests := map[string]struct {
		config string
		expSA  model.ServiceAccount
		expErr *regexp.Regexp
	}{
		"A correct configuration should create the service account.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
}
`,
			expSA: model.ServiceAccount{
				ID:    "ci-bot",
				Name:  "ci-bot",
				Token: "ops_fake_token_ci-bot",
			},
		},

		"A non-set name should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
}
`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty name should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = ""
}
`,
			expErr: regexp.MustCompile("Attribute name string length must be at least 1, got: 0"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare fake storage.
			path, deleteFn := getFakeRepoTmpFile("TestAccServiceAccountCreateDelete")
			defer deleteFn()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					assertServiceAccountOnFakeStorage(t, &test.expSA),
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "id", test.expSA.ID),
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "name", test.expSA.Name),
					// Token is sensitive; we can still check it exists in state via the helper.
				)
			}

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				CheckDestroy:             assertServiceAccountDeletedOnFakeStorage(t, test.expSA.ID),
				Steps: []resource.TestStep{
					{
						Config:      test.config,
						Check:       checks,
						ExpectError: test.expErr,
					},
				},
			})
		})
	}
}

// TestAccServiceAccountVaultAccess verifies that vault_access is round-tripped
// through state and forces replacement on change.
func TestAccServiceAccountVaultAccess(t *testing.T) {
	// Prepare fake storage.
	path, deleteFn := getFakeRepoTmpFile("TestAccServiceAccountVaultAccess")
	defer deleteFn()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	configCreate := `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"

  vault_access = {
    "vault-alpha" = ["read_items", "write_items"]
    "vault-beta"  = ["read_items"]
  }
}
`
	configReplace := `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"

  vault_access = {
    "vault-alpha" = ["read_items"]
  }
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "name", "ci-bot"),
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "vault_access.vault-alpha.#", "2"),
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "vault_access.vault-beta.#", "1"),
				),
			},
			{
				// Changing vault_access must force replacement (destroy + create).
				Config: configReplace,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "vault_access.vault-alpha.#", "1"),
					resource.TestCheckNoResourceAttr("onepasswordorg_service_account.test_sa", "vault_access.vault-beta"),
				),
			},
		},
	})
}

// TestAccServiceAccountVaultAccessValidation checks the config validator that
// enforces permission allowed-values, the write_items/share_items → read_items
// dependency, and the vault-identifier escape rule.
func TestAccServiceAccountVaultAccessValidation(t *testing.T) {
	tests := map[string]struct {
		config string
		expErr *regexp.Regexp
	}{
		"Unknown permission should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
  vault_access = {
    "vault-alpha" = ["bogus_perm"]
  }
}
`,
			expErr: regexp.MustCompile(`(?s)Invalid service account vault permission`),
		},
		"write_items without read_items should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
  vault_access = {
    "vault-alpha" = ["write_items"]
  }
}
`,
			expErr: regexp.MustCompile(`(?s)Missing required service account vault permission`),
		},
		"share_items without read_items should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
  vault_access = {
    "vault-alpha" = ["share_items"]
  }
}
`,
			expErr: regexp.MustCompile(`(?s)Missing required service account vault permission`),
		},
		"Vault identifier with ':' should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
  vault_access = {
    "bad:vault" = ["read_items"]
  }
}
`,
			expErr: regexp.MustCompile(`(?s)Invalid vault identifier`),
		},
		"Vault identifier with ',' should fail.": {
			config: `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
  vault_access = {
    "bad,vault" = ["read_items"]
  }
}
`,
			expErr: regexp.MustCompile(`(?s)Invalid vault identifier`),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			path, deleteFn := getFakeRepoTmpFile("TestAccServiceAccountVaultAccessValidation")
			defer deleteFn()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				Steps: []resource.TestStep{
					{
						Config:      test.config,
						ExpectError: test.expErr,
					},
				},
			})
		})
	}
}
func TestAccServiceAccountNameChangeReplace(t *testing.T) {
	// Prepare fake storage.
	path, deleteFn := getFakeRepoTmpFile("TestAccServiceAccountNameChangeReplace")
	defer deleteFn()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	configCreate := `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot"
}
`
	configUpdate := `
resource "onepasswordorg_service_account" "test_sa" {
  name = "ci-bot-renamed"
}
`

	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertServiceAccountOnFakeStorage(t, &model.ServiceAccount{
						ID:    "ci-bot",
						Name:  "ci-bot",
						Token: "ops_fake_token_ci-bot",
					}),
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "name", "ci-bot"),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertServiceAccountOnFakeStorage(t, &model.ServiceAccount{
						ID:    "ci-bot-renamed",
						Name:  "ci-bot-renamed",
						Token: "ops_fake_token_ci-bot-renamed",
					}),
					assertServiceAccountDeletedOnFakeStorage(t, "ci-bot"),
					resource.TestCheckResourceAttr("onepasswordorg_service_account.test_sa", "name", "ci-bot-renamed"),
				),
			},
		},
	})
}
