terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}
provider "onepasswordorg" {
  address = "mycompany.onepassword.com"
  email   = "bot+onepassword@mycompany.com"

  # Or use `OP_SECRET_KEY` env var.
  secret_key = var.op_secret_key
}
