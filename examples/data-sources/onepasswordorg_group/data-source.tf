data "onepasswordorg_group" "test" {
  name = "test-group"
}

output "group_test" {
  value = data.onepasswordorg_group.test
}
