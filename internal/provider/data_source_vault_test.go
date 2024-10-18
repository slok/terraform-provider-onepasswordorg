package provider_test

import (
	"context"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"
	"github.com/stretchr/testify/require"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccDataSourceVaultCorrect will check a vault can be used as data source.
func TestAccDataSourceVaultCorrect(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccDataSourceVaultCorrect")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
data "onepasswordorg_vault" "test" {
  name = "test-vault"
}
`
	// Prepare storage.
	repo := getFakeRepository(t)
	_, err := repo.CreateVault(context.TODO(), model.Vault{Name: "test-vault", Description: "Test vault"})
	require.NoError(t, err)
	defer func() { _ = repo.DeleteVault(context.TODO(), "test-vault") }()

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.onepasswordorg_vault.test", "id", "test-vault"), // Fake uses user name ID.
					resource.TestCheckResourceAttr("data.onepasswordorg_vault.test", "description", "Test vault"),
					resource.TestCheckResourceAttr("data.onepasswordorg_vault.test", "name", "test-vault"),
				),
			},
		},
	})
}

// TestAccDataSourceVaultMissing will check the datasource fails when the vault is missing.
func TestAccDataSourceVaultMissing(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccDataSourceVaultMissing")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
data "onepasswordorg_vault" "test" {
  name = "test-vault"
}
`

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("vault does not exists"),
			},
		},
	})
}
