terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}
provider "onepasswordorg" {
  address           = "my.onepassword.com"
  email             = "me@slok.dev"
  secret_key        = "test"
  fake_storage_path = "/tmp/tf-onepasswordorg-storage.json"

}

resource "onepasswordorg_user" "user0" {
  for_each = {
    "user0" : { name : "user0", email : "user00@slok.dev" }
    "user1" : { name : "user1", email : "user10@slok.dev" }
    "user2" : { name : "user2", email : "user20@slok.dev" }
  }

  name  = each.value.name
  email = each.value.email
}
