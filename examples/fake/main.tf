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

resource "onepasswordorg_user" "test" {
  for_each = {
    "user0" : { name : "User number 0", email : "user0@slok.dev" }
    "user1" : { name : "User number 1", email : "user1@slok.dev" }
    "user2" : { name : "User number 2", email : "user2@slok.dev" }
    "user3" : { name : "User number 3", email : "user3@slok.dev" }
  }

  name  = each.value.name
  email = each.value.email
}

resource "onepasswordorg_group" "test" {
  for_each = {
    "group0" : { name : "group-0", description : "Group zero" }
    "group1" : { name : "group-1", description : "Group one" }
    "group2" : { name : "group-2", description : "Group two" }
  }

  name        = each.value.name
  description = each.value.description
}


