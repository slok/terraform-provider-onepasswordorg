package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccGroupCreateDelete will check a group is created and deleted.
func TestAccGroupCreateDelete(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccGroupCreateDelete")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
resource "onepasswordorg_group" "test_group" {
  name  = "test-group"
  description = "Test group"
}
`
	// Fake repo group IDs are based on name.
	expGroup := model.Group{
		ID:          "test-group",
		Name:        "test-group",
		Description: "Test group",
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             assertGroupDeletedOnFakeStorage(t, "test-group"), // Fake uses name as IDs.
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertGroupOnFakeStorage(t, &expGroup),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "id", "test-group"), // Fake uses name as IDs.
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "name", "test-group"),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "description", "Test group"),
				),
			},
		},
	})
}

// TestAcc_GroupUpdateDescription will check a group can update its description after its creation.
func TestAccGroupUpdateDescription(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccGroupUpdateName")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	configCreate := `
resource "onepasswordorg_group" "test_group" {
  name  	  = "test-group"
  description = "Test group"
}
`
	configUpdate := `
resource "onepasswordorg_group" "test_group" {
  name  	  = "test-group"
  description = "Test group modified"
}
`

	// Fake repo IDs are based on emails.
	expGroupCreate := model.Group{
		ID:          "test-group",
		Name:        "test-group",
		Description: "Test group",
	}

	expGroupUpdate := model.Group{
		ID:          "test-group",
		Name:        "test-group",
		Description: "Test group modified",
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertGroupOnFakeStorage(t, &expGroupCreate),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertGroupOnFakeStorage(t, &expGroupUpdate),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "id", "test-group"), // Fake uses name as IDs.
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "name", "test-group"),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "description", "Test group modified"),
				),
			},
		},
	})
}
