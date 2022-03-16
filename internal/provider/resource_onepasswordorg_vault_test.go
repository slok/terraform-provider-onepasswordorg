package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccVaultCreateDelete will check a vault is created and deleted.
func TestAccVaultCreateDelete(t *testing.T) {
	tests := map[string]struct {
		config   string
		expVault model.Vault
		expErr   *regexp.Regexp
	}{
		"A correct configuration should execute correctly.": {
			config: `
resource "onepasswordorg_vault" "test" {
  name  = "test-vault"
  description = "Test vault"
}
`,
			expVault: model.Vault{
				ID:          "test-vault",
				Name:        "test-vault",
				Description: "Test vault",
			},
		},

		"Description should fallback if not set.": {
			config: `
resource "onepasswordorg_vault" "test" {
  name  = "test-vault"
 
}
		`,
			expVault: model.Vault{
				ID:          "test-vault",
				Name:        "test-vault",
				Description: "Managed by Terraform",
			},
		},

		"A non set name should fail.": {
			config: `
resource "onepasswordorg_vault" "test" {
	description = "Test vault"
}
		`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty name should fail.": {
			config: `
resource "onepasswordorg_vault" "test" {
	description = "Test vault"
	name = ""
}
		`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare fake storage.
			path, delete := getFakeRepoTmpFile("TestAccVaultCreateDelete")
			defer delete()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					assertVaultOnFakeStorage(t, &test.expVault),
					resource.TestCheckResourceAttr("onepasswordorg_vault.test", "id", test.expVault.ID),
					resource.TestCheckResourceAttr("onepasswordorg_vault.test", "name", test.expVault.Name),
					resource.TestCheckResourceAttr("onepasswordorg_vault.test", "description", test.expVault.Description),
				)
			}

			// Execute test.
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				CheckDestroy:             assertVaultDeletedOnFakeStorage(t, test.expVault.ID),
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

// TestAccVaultUpdateDescription will check a vault can update its description after its creation.
func TestAccVaultUpdateDescription(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccVaultUpdateName")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	configCreate := `
resource "onepasswordorg_vault" "test" {
  name  	  = "test-vault"
  description = "Test vault"
}
`
	configUpdate := `
resource "onepasswordorg_vault" "test" {
  name  	  = "test-vault"
  description = "Test vault modified"
}
`

	// Fake repo IDs are based on emails.
	expVaultCreate := model.Vault{
		ID:          "test-vault",
		Name:        "test-vault",
		Description: "Test vault",
	}

	expVaultUpdate := model.Vault{
		ID:          "test-vault",
		Name:        "test-vault",
		Description: "Test vault modified",
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertVaultOnFakeStorage(t, &expVaultCreate),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertVaultOnFakeStorage(t, &expVaultUpdate),
					resource.TestCheckResourceAttr("onepasswordorg_vault.test", "id", "test-vault"), // Fake uses name as IDs.
					resource.TestCheckResourceAttr("onepasswordorg_vault.test", "name", "test-vault"),
					resource.TestCheckResourceAttr("onepasswordorg_vault.test", "description", "Test vault modified"),
				),
			},
		},
	})
}
