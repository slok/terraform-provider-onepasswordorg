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
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccGroupMemberCreateDelete")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
  role     = "member"
}
`
	// Fake repo IDs are based on group id + user id.
	expMember := model.Membership{
		GroupID: "test-group-id",
		UserID:  "test-user-id",
		Role:    model.MembershipRoleMember,
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             assertGroupMemberDeletedOnFakeStorage(t, "test-group-id", "test-user-id"), // Fake uses group and user IDs.
		Steps: []resource.TestStep{
			{
				Config: config,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertGroupMemberOnFakeStorage(t, &expMember),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "id", "test-group-id/test-user-id"), // Fake uses group and user IDs.
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "group_id", "test-group-id"),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "user_id", "test-user-id"),
					resource.TestCheckResourceAttr("onepasswordorg_group_member.test_member", "role", "member"),
				),
			},
		},
	})
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
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
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

// TestAccGroupMemberInvalidRole will check a group member doesn't allow invalid roles.
func TestAccGroupMemberInvalidRole(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccGroupMemberInvalidRole")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.
	config := `
resource "onepasswordorg_group_member" "test_member" {
  group_id = "test-group-id"
  user_id  = "test-user-id"
  role     = "invalid-role"
}
`
	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		CheckDestroy:             assertGroupMemberDeletedOnFakeStorage(t, "test-group-id", "test-user-id"), // Fake uses group and user IDs.
		Steps: []resource.TestStep{
			{
				Config:      config,
				ExpectError: regexp.MustCompile(`the role "[^"]*" is invalid`),
			},
		},
	})
}
