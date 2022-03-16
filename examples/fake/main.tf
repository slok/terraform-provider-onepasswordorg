terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}

provider "onepasswordorg" {
  fake_storage_path = "/tmp/tf-onepasswordorg-storage.json"
}

locals {
  users = {
    "user0" : { name : "User number 0", email : "user0@slok.dev" },
    "user1" : { name : "User number 1", email : "user1@slok.dev" },
    "user2" : { name : "User number 2", email : "user2@slok.dev" },
    "user3" : { name : "User number 3", email : "user3@slok.dev" },
    "user4" : { name : "User number 4", email : "user4@slok.dev" },
    "user5" : { name : "User number 5", email : "user5@slok.dev" },
    "user6" : { name : "User number 6", email : "user6@slok.dev" },
    "user7" : { name : "User number 7", email : "user7@slok.dev" },
    "user8" : { name : "User number 8", email : "user8@slok.dev" },
    "user9" : { name : "User number 9", email : "user9@slok.dev" },
  }

  groups = {
    "group0" : { name : "group-0", description : "Group zero" },
    "group1" : { name : "group-1", description : "Group one" },
    "group2" : { name : "group-2", description : null },
  }

  vaults = {
    "vault0" : { name : "vault-0", description : "Vault zero" },
    "vault1" : { name : "vault-1", description : "Vault one" },
    "vault2" : { name : "vault-2", description : null },
    "vault3" : { name : "vault-3", description : "Vault three" },
    "vault4" : { name : "vault-4", description : "Vault four" },
    "vault5" : { name : "vault-5", description : "Vault five5" },
  }

  members = {
    "group0-user0" : { user_id : "user0", group_id : "group0", role : "member" },
    "group0-user1" : { user_id : "user1", group_id : "group0", role : "member" },
    "group1-user0" : { user_id : "user0", group_id : "group1", role : null },
    "group1-user2" : { user_id : "user2", group_id : "group1", role : "manager" },
  }
}

# Users.
resource "onepasswordorg_user" "test" {
  for_each = local.users

  name  = each.value.name
  email = each.value.email
}

# Groups.
resource "onepasswordorg_group" "test" {
  for_each = local.groups

  name        = each.value.name
  description = each.value.description
}

resource "onepasswordorg_group_member" "test" {
  for_each = local.members

  group_id = each.value.group_id
  user_id  = each.value.user_id
  role     = each.value.role
}

# Vaults.
resource "onepasswordorg_vault" "test" {
  for_each = local.vaults

  name        = each.value.name
  description = each.value.description
}


# Data.
data "onepasswordorg_user" "user4" {
  email = onepasswordorg_user.test["user4"].email
}

output "user4" {
  value = data.onepasswordorg_user.user4
}


data "onepasswordorg_group" "group2" {
  name = onepasswordorg_group.test["group2"].name
}

output "group2" {
  value = data.onepasswordorg_group.group2
}
