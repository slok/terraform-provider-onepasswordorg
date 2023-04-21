package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccUserCreateDelete will check a user is created and deleted.
func TestAccUserCreateDelete(t *testing.T) {
	tests := map[string]struct {
		config  string
		expUser model.User
		expErr  *regexp.Regexp
	}{
		"A correct configuration should execute correctly.": {
			config: `
resource "onepasswordorg_user" "test_user" {
  name  = "Test user"
  email = "testuser@test.test"
}
`,
			expUser: model.User{
				ID:    "testuser@test.test",
				Name:  "Test user",
				Email: "testuser@test.test",
			},
		},

		"A non set name should fail.": {
			config: `
resource "onepasswordorg_user" "test_user" {
  email = "testuser@test.test"
}
`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty user should fail.": {
			config: `
resource "onepasswordorg_user" "test_user" {
  name  = ""
  email = "testuser@test.test"
}
`,
			expErr: regexp.MustCompile("Error: expected \"name\" to not be an empty string"),
		},

		"A non set email should fail.": {
			config: `
resource "onepasswordorg_user" "test_user" {
  name  = "Test user"
}
`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty email should fail.": {
			config: `
resource "onepasswordorg_user" "test_user" {
  name  = "Test user"
  email = ""
}
`,
			expErr: regexp.MustCompile("Error: expected \"email\" to not be an empty string"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare fake storage.
			path, delete := getFakeRepoTmpFile("TestAccUserCreateDelete")
			defer delete()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					assertUserOnFakeStorage(t, &test.expUser),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "id", test.expUser.ID),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "name", test.expUser.Name),
					resource.TestCheckResourceAttr("onepasswordorg_user.test_user", "email", test.expUser.Email),
				)
			}

			// Execute test.
			resource.Test(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: assertUserDeletedOnFakeStorage(t, test.expUser.Email),
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
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
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
