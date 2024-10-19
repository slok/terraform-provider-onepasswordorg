package provider_test

import (
	"fmt"
	"os"
	"regexp"
	"testing"

	"github.com/hashicorp/terraform-plugin-testing/helper/resource"

	"github.com/slok/terraform-provider-onepasswordorg/internal/model"
	"github.com/slok/terraform-provider-onepasswordorg/internal/provider"
)

// TestAccVaultGroupAccessDelete will check a vault group access is created and deleted.
func TestAccVaultGroupAccessCreateDelete(t *testing.T) {
	tests := map[string]struct {
		config string
		expID  string
		expVGA model.VaultGroupAccess
		expErr *regexp.Regexp
	}{
		"A correct configuration should execute correctly.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	  allow_viewing = true
	  allow_editing = true
	  allow_managing = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID: "test-vault-id",
				GroupID: "test-group-id",
				Permissions: model.AccessPermissions{
					AllowViewing:  true,
					AllowEditing:  true,
					AllowManaging: true,
				},
			},
		},

		"Permission allow_viewing check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	allow_viewing = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{AllowViewing: true},
			},
		},

		"Permission allow_editing check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	allow_editing = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{AllowEditing: true},
			},
		},

		"Permission allow_managing check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	allow_managing = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{AllowManaging: true},
			},
		},

		"Permission view_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	view_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ViewItems: true},
			},
		},

		"Permission create_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	create_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{CreateItems: true},
			},
		},

		"Permission edit_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	edit_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{EditItems: true},
			},
		},

		"Permission archive_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	archive_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ArchiveItems: true},
			},
		},

		"Permission delete_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	delete_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{DeleteItems: true},
			},
		},

		"Permission view_and_copy_passwords check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	view_and_copy_passwords = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ViewAndCopyPasswords: true},
			},
		},

		"Permission view_item_history check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	view_item_history = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ViewItemHistory: true},
			},
		},

		"Permission import_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	import_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ImportItems: true},
			},
		},

		"Permission export_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	export_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ExportItems: true},
			},
		},

		"Permission copy_and_share_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	copy_and_share_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{CopyAndShareItems: true},
			},
		},

		"Permission print_items check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	print_items = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{PrintItems: true},
			},
		},

		"Permission manage_vault check.": {
			config: `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	manage_vault = true
  }
}
`,
			expID: "test-vault-id/test-group-id",
			expVGA: model.VaultGroupAccess{
				VaultID:     "test-vault-id",
				GroupID:     "test-group-id",
				Permissions: model.AccessPermissions{ManageVault: true},
			},
		},
	}

	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			// Prepare fake storage.
			path, delete := getFakeRepoTmpFile("TestAccVaultGroupAccessCreateDelete")
			defer delete()
			_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

			// Prepare non error checks.
			var checks resource.TestCheckFunc
			if test.expErr == nil {
				checks = resource.ComposeAggregateTestCheckFunc(
					assertVaultGroupAccessOnFakeStorage(t, &test.expVGA),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "id", test.expID),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "vault_id", test.expVGA.VaultID),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "group_id", test.expVGA.GroupID),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.allow_viewing", fmt.Sprintf("%t", test.expVGA.Permissions.AllowViewing)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.allow_editing", fmt.Sprintf("%t", test.expVGA.Permissions.AllowEditing)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.allow_managing", fmt.Sprintf("%t", test.expVGA.Permissions.AllowManaging)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.view_items", fmt.Sprintf("%t", test.expVGA.Permissions.ViewItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.create_items", fmt.Sprintf("%t", test.expVGA.Permissions.CreateItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.edit_items", fmt.Sprintf("%t", test.expVGA.Permissions.EditItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.archive_items", fmt.Sprintf("%t", test.expVGA.Permissions.ArchiveItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.delete_items", fmt.Sprintf("%t", test.expVGA.Permissions.DeleteItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.view_and_copy_passwords", fmt.Sprintf("%t", test.expVGA.Permissions.ViewAndCopyPasswords)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.view_item_history", fmt.Sprintf("%t", test.expVGA.Permissions.ViewItemHistory)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.import_items", fmt.Sprintf("%t", test.expVGA.Permissions.ImportItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.export_items", fmt.Sprintf("%t", test.expVGA.Permissions.ExportItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.copy_and_share_items", fmt.Sprintf("%t", test.expVGA.Permissions.CopyAndShareItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.print_items", fmt.Sprintf("%t", test.expVGA.Permissions.PrintItems)),
					resource.TestCheckResourceAttr("onepasswordorg_vault_group_access.test", "permissions.manage_vault", fmt.Sprintf("%t", test.expVGA.Permissions.ManageVault)),
				)
			}

			// Execute test.
			resource.Test(t, resource.TestCase{
				PreCheck:                 func() { testAccPreCheck(t) },
				ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
				CheckDestroy:             assertVaultGroupAccessDeletedOnFakeStorage(t, test.expVGA.VaultID, test.expVGA.GroupID),
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

// TestAccVaultGroupAccessrUpdateRole will check a membership is can update the role.
func TestAccVaultGroupAccessrUpdateRole(t *testing.T) {
	// Prepare fake storage.
	path, delete := getFakeRepoTmpFile("TestAccVaultGroupAccessrUpdateRole")
	defer delete()
	_ = os.Setenv(provider.EnvVarOpFakeStoragePath, path)

	// Test tf data.

	configCreate := `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	  allow_viewing = true
	  allow_editing = true
	  allow_managing = true
  }
}`
	configUpdate := `
resource "onepasswordorg_vault_group_access" "test" {
  vault_id  = "test-vault-id"
  group_id = "test-group-id" 
  permissions = {
	  allow_viewing = false
	  allow_editing = true
	  view_and_copy_passwords = true
	  print_items = true
  }
}`

	expVGACreate := model.VaultGroupAccess{
		VaultID: "test-vault-id",
		GroupID: "test-group-id",
		Permissions: model.AccessPermissions{
			AllowViewing:  true,
			AllowEditing:  true,
			AllowManaging: true,
		},
	}

	expVGAUpdate := model.VaultGroupAccess{
		VaultID: "test-vault-id",
		GroupID: "test-group-id",
		Permissions: model.AccessPermissions{
			AllowEditing:         true,
			ViewAndCopyPasswords: true,
			PrintItems:           true,
		},
	}

	// Execute test.
	resource.Test(t, resource.TestCase{
		PreCheck:                 func() { testAccPreCheck(t) },
		ProtoV6ProviderFactories: testAccProtoV6ProviderFactories,
		Steps: []resource.TestStep{
			{
				Config: configCreate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertVaultGroupAccessOnFakeStorage(t, &expVGACreate),
				),
			},
			{
				Config: configUpdate,
				Check: resource.ComposeAggregateTestCheckFunc(
					assertVaultGroupAccessOnFakeStorage(t, &expVGAUpdate),
				),
			},
		},
	})
}
