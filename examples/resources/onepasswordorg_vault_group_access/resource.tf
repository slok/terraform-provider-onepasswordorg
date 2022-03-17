resource "onepasswordorg_vault" "vault0" {
  name        = "vault-0"
  description = "Vault 0"
}

resource "onepasswordorg_group" "group0" {
  name        = "group-0"
  description = "Group 0"
}

resource "onepasswordorg_vault_group_access" "business_full" {
  user_id  = onepasswordorg_vault.vault0.id
  group_id = onepasswordorg_group.group0.id
  permissions = {
    view_items              = true
    create_items            = true
    edit_items              = true
    archive_items           = true
    delete_items            = true
    view_and_copy_passwords = true
    view_item_history       = true
    import_items            = true
    export_items            = true
    copy_and_share_items    = true
    print_items             = true
    manage_vault            = true
  }
}

resource "onepasswordorg_group" "group1" {
  name        = "group-1"
  description = "Group 1"
}

resource "onepasswordorg_vault_group_access" "team_view" {
  user_id  = onepasswordorg_vault.vault0.id
  group_id = onepasswordorg_group.group1.id
  permissions = {
    allow_viewing = true
  }
}

resource "onepasswordorg_group" "group2" {
  name        = "group-2"
  description = "Group 2"
}

resource "onepasswordorg_vault_group_access" "business_regular" {
  user_id  = onepasswordorg_vault.vault0.id
  group_id = onepasswordorg_group.group2.id
  permissions = {
    view_items              = true
    create_items            = true
    edit_items              = true
    archive_items           = true
    delete_items            = true
    view_and_copy_passwords = true
    view_item_history       = true
    import_items            = true
    export_items            = true
    copy_and_share_items    = true
    print_items             = true
  }
}


resource "onepasswordorg_group" "group3" {
  name        = "group-3"
  description = "Group 3"
}

resource "onepasswordorg_vault_group_access" "business_manage" {
  user_id  = onepasswordorg_vault.vault0.id
  group_id = onepasswordorg_group.group3.id
  permissions = {
    manage_vault = true
  }
}

