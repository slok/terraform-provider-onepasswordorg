terraform {
  required_providers {
    onepasswordorg = {
      source = "slok/onepasswordorg"
    }
  }
}

provider "onepasswordorg" {
  address = "mycompany.1password.com"
  email   = "bot+onepassword@mycompany.com"

  # Or use `OP_SECRET_KEY` env var.
  secret_key = var.op_secret_key

  # Or use `OP_PASSWORD` env var.
  password = var.op_password
}

resource "onepasswordorg_user" "test" {
  name  = "A test account"
  email = "test@mycompany.com"
}

resource "onepasswordorg_group" "test" {
  name        = "test-tf"
  description = "Group for testing terraform cloud"
}

resource "onepasswordorg_group_member" "test_test" {
  group_id = onepasswordorg_group.test.id
  user_id  = onepasswordorg_user.test.id
  role     = "member"
}
