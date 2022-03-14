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
  }

  groups = {
    "group0" : { name : "group-0", description : "Group zero" },
    "group1" : { name : "group-1", description : "Group one" },
    "group2" : { name : "group-2", description : "Group two" },
  }

  members = {
    "group0-user0" : { user_id : "user0", group_id : "group0", role : "member" },
    "group0-user1" : { user_id : "user1", group_id : "group0", role : "member" },
    "group1-user0" : { user_id : "user0", group_id : "group1", role : "manager" },
    "group1-user2" : { user_id : "user2", group_id : "group1", role : "manager" },
  }
}

resource "onepasswordorg_user" "test" {
  for_each = local.users

  name  = each.value.name
  email = each.value.email
}

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
