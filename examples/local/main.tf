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
