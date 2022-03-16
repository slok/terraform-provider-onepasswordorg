package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccGroupCreateDelete will check a group is created and deleted.
func TestAccGroupCreateDelete(t *testing.T) {
	tests := map[string]struct {
		config   string
		expGroup model.Group
		expErr   *regexp.Regexp
	}{
		"A correct configuration should execute correctly.": {
			config: `
resource "onepasswordorg_group" "test_group" {
  name  = "test-group"
  description = "Test group"
}
`,
			expGroup: model.Group{
				ID:          "test-group",
				Name:        "test-group",
				Description: "Test group",
			},
		},

		"Description should fallback if not set.": {
			config: `
resource "onepasswordorg_group" "test_group" {
  name  = "test-group"
}
`,
			expGroup: model.Group{
				ID:          "test-group",
				Name:        "test-group",
				Description: "Managed by Terraform",
			},
		},

		"A non set name should fail.": {
			config: `
resource "onepasswordorg_group" "test_group" {
  description = "Test group"
}
`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty name should fail.": {
			config: `
resource "onepasswordorg_group" "test_group" {
  description = "Test group"
  name = ""
}
`,
			expErr: regexp.MustCompile("Attribute can't be empty"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare fake storage.
			path, delete := getFakeRepoTmpFile("TestAccGroupCreateDelete")
			defer delete()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					assertGroupOnFakeStorage(t, &test.expGroup),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "id", test.expGroup.ID),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "name", test.expGroup.Name),
					resource.TestCheckResourceAttr("onepasswordorg_group.test_group", "description", test.expGroup.Description),
				)
			}

			// Execute test.
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				CheckDestroy:             assertGroupDeletedOnFakeStorage(t, test.expGroup.ID),
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

// TestAccGroupUpdateDescription will check a group can update its description after its creation.
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
