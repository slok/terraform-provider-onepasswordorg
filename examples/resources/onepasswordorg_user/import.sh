# Get user UUID with op CLI.
op get user user0@slok.dev

# Import.
terraform import onepasswordorg_user.user0 ${ONEPASSWORD_UUID}
