resource "onepasswordorg_user" "user0" {
  name  = "User zero"
  email = "user0@slok.dev"
}

resource "onepasswordorg_group" "test_group" {
  name        = "test-group"
  description = "Group for testing"
}

resource "onepasswordorg_group_member" "test_group_user0" {
  group_id = onepasswordorg_group.test_group.id
  user_id  = onepasswordorg_user.user0.id
  role     = "member"
}
