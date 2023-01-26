package provider_test

import (
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccGroupMemberCreateDelete will check a group member is created and deleted.
func TestAccGroupMemberCreateDelete(t *testing.T) {
	tests := map[string]struct {
		config    string
		expID     string
		expRole   string
		expMember model.Membership
		expErr    *regexp.Regexp
	}{
		"A correct configuration should execute correctly.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
  role     = "manager"
}
`,
			expID:   "test-group-id/test-user-id",
			expRole: "manager",
			expMember: model.Membership{
				GroupID: "test-group-id",
				UserID:  "test-user-id",
				Role:    model.MembershipRoleManager,
			},
		},

		"If role not set it should fallback to a default role.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
}
`,
			expID:   "test-group-id/test-user-id",
			expRole: "member",
			expMember: model.Membership{
				GroupID: "test-group-id",
				UserID:  "test-user-id",
				Role:    model.MembershipRoleMember,
			},
		},

		"An invalid role should fail.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
  role     = "invalid-role"
}
`,
			expErr: regexp.MustCompile(`the role "invalid-role" is invalid`),
		},

		"A non set group id should fail.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  user_id  = "test-user-id"
  role     = "invalid-role"
}
`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty group id should fail.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  group_id = ""
  user_id  = "test-user-id"
  role     = "invalid-role"
}
`,
			expErr: regexp.MustCompile("Error: expected \"group_id\" to not be an empty string"),
		},

		"A non set user id should fail.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  role     = "invalid-role"
}
`,
			expErr: regexp.MustCompile("Missing required argument"),
		},

		"An empty user id should fail.": {
			config: `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = ""
  role     = "invalid-role"
}
`,
			expErr: regexp.MustCompile("Error: expected \"user_id\" to not be an empty string"),
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare fake storage.
			path, delete := getFakeRepoTmpFile("TestAccGroupMemberCreateDelete")
			defer delete()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					assertGroupMemberOnFakeStorage(t, &test.expMember),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "id", test.expID),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "group_id", test.expMember.GroupID),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "user_id", test.expMember.UserID),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "role", test.expRole),
				)
			}

			// Execute test.
			resource.Test(t, resource.TestCase{
				PreCheck:     func() { testAccPreCheck(t) },
				Providers:    testAccProviders,
				CheckDestroy: assertGroupMemberDeletedOnFakeStorage(t, test.expMember.GroupID, test.expMember.UserID),
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

// TestAccGroupMemberUpdateRole will check a membership is can update the role.
func TestAccGroupMemberUpdateRole(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccGroupMemberUpdateRole")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	configCreate := `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
  role     = "member"
}
`
	configUpdate := `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
  role     = "manager"
}
`
	// Fake repo IDs are based on group id + user id.
	expMemberCreate := model.Membership{
		GroupID: "test-group-id",
		UserID:  "test-user-id",
		Role:    model.MembershipRoleMember,
	}

	expMemberUpdate := model.Membership{
		GroupID: "test-group-id",
		UserID:  "test-user-id",
		Role:    model.MembershipRoleManager,
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertGroupMemberOnFakeStorage(t, &expMemberCreate),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertGroupMemberOnFakeStorage(t, &expMemberUpdate),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "id", "test-group-id/test-user-id"), // Fake uses group and user IDs.
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "group_id", "test-group-id"),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "user_id", "test-user-id"),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "role", "manager"),
				),
			},
		},
	})
}
