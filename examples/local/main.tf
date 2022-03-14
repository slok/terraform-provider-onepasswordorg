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
    "infrastructure-test" : { name : "1password test", email : "infrastructure+test@fonoa.com" }
    "infrastructure-test2" : { name : "1password test2", email : "infrastructure+test2@fonoa.com" }
  }

  name  = each.value.name
  email = each.value.email
}

resource "onepasswordorg_group" "test_tf" {
  name        = "test-tf"
  description = "TF group test"
}
