package provider_test

import (
	"os"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccUserCreateDelete will check a user is created and deleted.
func TestAccUserCreateDelete(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccUserCreateDelete")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
resource "onepasswordorg_user" "test_user" {
  name  = "Test user"
  email = "testuser@test.test"
}
`
	// Fake repo IDs are based on emails.
	expUser := model.User{
		ID:    "testuser@test.test",
		Name:  "Test user",
		Email: "testuser@test.test",
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             assertUserDeletedOnFakeStorage(t, "testuser@test.test"), // Fake uses email as IDs.
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertUserOnFakeStorage(t, &expUser),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "id", "testuser@test.test"), // Fake uses email as IDs.
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "name", "Test user"),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "email", "testuser@test.test"),
				),
			},
		},
	})
}

// TestAccUserUpdateName will check a user can update its name after its creation.
func TestAccUserUpdateName(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccUserUpdateName")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	configCreate := `
resource "onepasswordorg_user" "test_user" {
  name  = "Test user"
  email = "testuser@test.test"
}
`
	configUpdate := `
resource "onepasswordorg_user" "test_user" {
  name  = "Test user modified"
  email = "testuser@test.test"
}
`

	// Fake repo IDs are based on emails.
	expUserCreate := model.User{
		ID:    "testuser@test.test",
		Name:  "Test user",
		Email: "testuser@test.test",
	}

	expUserUpdate := model.User{
		ID:    "testuser@test.test",
		Name:  "Test user modified",
		Email: "testuser@test.test",
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertUserOnFakeStorage(t, &expUserCreate),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertUserOnFakeStorage(t, &expUserUpdate),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "id", "testuser@test.test"), // Fake uses email as IDs.
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "name", "Test user modified"),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "email", "testuser@test.test"),
				),
			},
		},
	})
}
