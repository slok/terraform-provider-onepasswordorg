# Go to the website and get the UUID from the URL or use the `op` cli:
op get group test-group

# Import.
terraform import onepasswordorg_group.group0 ${ONEPASSWORD_UUID}
