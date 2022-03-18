terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}

provider "onepasswordorg" {}

# Resources.
resource "onepasswordorg_user" "test" {
  name  = "1password test 3"
  email = "infrastructure+test3@fonoa.com"
}

resource "onepasswordorg_group" "test" {
  name        = "test-tf2"
  description = "TF group test!!"
}

resource "onepasswordorg_group_member" "test" {
  user_id  = onepasswordorg_user.test.id
  group_id = onepasswordorg_group.test.id
  role     = "manager"
}

resource "onepasswordorg_vault" "test" {
  name        = "test-tf"
  description = "Terraform test vault"
}

resource "onepasswordorg_vault_group_access" "test" {
  vault_id = onepasswordorg_vault.test.id
  group_id = onepasswordorg_group.test.id
  permissions = {
    archive_items           = true
    copy_and_share_items    = true
    create_items            = true
    delete_items            = true
    edit_items              = true
    export_items            = true
    import_items            = true
    manage_vault            = false
    print_items             = true
    view_and_copy_passwords = true
    view_item_history       = true
    view_items              = true
  }
}

resource "onepasswordorg_vault_user_access" "test" {
  vault_id = onepasswordorg_vault.test.id
  user_id  = onepasswordorg_user.test.id
  permissions = {
    archive_items           = false
    copy_and_share_items    = false
    create_items            = false
    delete_items            = false
    edit_items              = false
    export_items            = false
    import_items            = false
    manage_vault            = false
    print_items             = false
    view_and_copy_passwords = true
    view_item_history       = true
    view_items              = true
  }
}

# Data.
data "onepasswordorg_user" "test" {
  email = onepasswordorg_user.test.email
}

output "user_test" {
  value = data.onepasswordorg_user.test
}

data "onepasswordorg_group" "test" {
  name = onepasswordorg_group.test.name
}

output "group_test" {
  value = data.onepasswordorg_group.test
}

data "onepasswordorg_vault" "test" {
  name = onepasswordorg_vault.test.name
}

output "vault_test" {
  value = data.onepasswordorg_vault.test
}
