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

// TestAccDataSourceUserCorrect will check a user can be used as data source.
func TestAccDataSourceUserCorrect(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccDataSourceUserCorrect")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
data "onepasswordorg_user" "test" {
  email = "test@slok.dev"
}
`
	// Prepare storage.
	repo := getFakeRepository(t)
	_, err := repo.CreateUser(context.TODO(), model.User{Email: "test@slok.dev", Name: "Test user"})
	require.NoError(t, err)
	defer func() { _ = repo.DeleteUser(context.TODO(), "test@slok.dev") }()

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.onepasswordorg_user.test", "id", "test@slok.dev"), // Fake uses user email ID.
					resource.TestCheckResourceAttr("data.onepasswordorg_user.test", "email", "test@slok.dev"),
					resource.TestCheckResourceAttr("data.onepasswordorg_user.test", "name", "Test user"),
				),
			},
		},
	})
}

// TestAccDataSourceUserMissing will check the datasource fails when the user is missing.
func TestAccDataSourceUserMissing(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccDataSourceUserMissing")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
data "onepasswordorg_user" "test" {
  email = "test@slok.dev"
}
`

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("user does not exists"),
			},
		},
	})
}
