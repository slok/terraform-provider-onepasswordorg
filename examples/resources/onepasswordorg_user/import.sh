# Go to the website and get the UUID from the URL or use the `op` cli:
op user get user0@slok.dev

# Import.
terraform import onepasswordorg_user.user0 ${ONEPASSWORD_UUID}
