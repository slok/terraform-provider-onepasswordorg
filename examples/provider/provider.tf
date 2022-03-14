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
