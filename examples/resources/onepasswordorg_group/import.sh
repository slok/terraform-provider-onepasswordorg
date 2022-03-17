# Go to the website and get the UUID from the URL or use the `op` cli:
op group get test-group

# Import.
terraform import onepasswordorg_group.group0 ${ONEPASSWORD_UUID}
