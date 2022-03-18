resource "onepasswordorg_vault" "vault0" {
  name        = "vault-0"
  description = "Vault 0"
}

resource "onepasswordorg_user" "user0" {
  name  = "user-0"
  email = "user0@slok.dev"
}

resource "onepasswordorg_vault_user_access" "business_full" {
  vault_id = onepasswordorg_vault.vault0.id
  user_id  = onepasswordorg_user.user0.id
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

resource "onepasswordorg_user" "user1" {
  name  = "user-1"
  email = "user1@slok.dev"
}

resource "onepasswordorg_vault_user_access" "team_view" {
  vault_id = onepasswordorg_vault.vault0.id
  user_id  = onepasswordorg_user.user1.id
  permissions = {
    allow_viewing = true
  }
}

resource "onepasswordorg_user" "user2" {
  name  = "user-2"
  email = "user2@slok.dev"
}

resource "onepasswordorg_vault_user_access" "business_regular" {
  vault_id = onepasswordorg_vault.vault0.id
  user_id  = onepasswordorg_user.user2.id
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


resource "onepasswordorg_user" "user3" {
  name  = "user-3"
  email = "user3@slok.dev"
}

resource "onepasswordorg_vault_user_access" "business_manage" {
  vault_id = onepasswordorg_vault.vault0.id
  user_id  = onepasswordorg_user.user3.id
  permissions = {
    manage_vault = true
  }
}
