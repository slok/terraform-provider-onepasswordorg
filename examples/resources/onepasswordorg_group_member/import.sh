# You will need the group ID and the user ID (<group id>/<user id>).
#
# Go to the website and get the UUID from the URL or use the `op` cli:
op get group test-group
op get user user0@slok.dev

# Import.
terraform import onepasswordorg_group_member.group0_member0 ${OP_GROUP_UUID}/${OP_USER_UUID}
