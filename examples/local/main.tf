terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}
provider "onepasswordorg" {}

resource "onepasswordorg_user" "test" {
  for_each = {
    "infrastructure-test" : { name : "1password test 2", email : "infrastructure+test@fonoa.com" }
  }

  name  = each.value.name
  email = each.value.email
}

resource "onepasswordorg_group" "test_tf" {
  name        = "test-tf"
  description = "TF group test"
}

resource "onepasswordorg_group_member" "test" {
  user_id  = onepasswordorg_user.test["infrastructure-test"].id
  group_id = onepasswordorg_group.test_tf.id
  role     = "manager"
}
