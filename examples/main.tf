terraform {
  required_providers {
    onepasswordorg = {
      source = "slok.dev/tf/onepasswordorg"
    }
  }
}
provider "onepasswordorg" {
  address    = "my.onepassword.com"
  email      = "me@slok.dev"
  secret_key = "test"
}

resource "onepasswordorg_user" "pepe_honka" {
  email = "pepe-honka@fonoa.com"
}
