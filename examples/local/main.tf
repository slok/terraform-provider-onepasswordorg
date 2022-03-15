terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}

provider "onepasswordorg" {}

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

data "onepasswordorg_user" "test" {
  email = onepasswordorg_user.test.email
}

output "user_test_id" {
  value = data.onepasswordorg_user.test
}
