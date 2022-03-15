data "onepasswordorg_user" "test" {
  email = "user0@slok.dev"
}

output "user_test_id" {
  value = data.onepasswordorg_user.test
}
