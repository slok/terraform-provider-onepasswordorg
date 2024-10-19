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

// TestAccDataSourceGroupCorrect will check a group can be used as data source.
func TestAccDataSourceGroupCorrect(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccDataSourceGroupCorrect")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
data "onepasswordorg_group" "test" {
  name = "test-group"
}
`
	// Prepare storage.
	repo := getFakeRepository(t)
	_, err := repo.CreateGroup(context.TODO(), model.Group{Name: "test-group", Description: "Test group"})
	require.NoError(t, err)
	defer func() { _ = repo.DeleteGroup(context.TODO(), "test-group") }()

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					resource.TestCheckResourceAttr("data.onepasswordorg_group.test", "id", "test-group"), // Fake uses user name ID.
					resource.TestCheckResourceAttr("data.onepasswordorg_group.test", "description", "Test group"),
					resource.TestCheckResourceAttr("data.onepasswordorg_group.test", "name", "test-group"),
				),
			},
		},
	})
}

// TestAccDataSourceGroupMissing will check the datasource fails when the group is missing.
func TestAccDataSourceGroupMissing(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccDataSourceGroupMissing")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
data "onepasswordorg_group" "test" {
  name = "test-group"
}
`

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile("group does not exists"),
			},
		},
	})
}
