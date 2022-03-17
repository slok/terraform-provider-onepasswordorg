# Go to the website and get the UUID from the URL or use the `op` cli:
op vault get test-vault

# Import.
terraform import onepasswordorg_vault.vault0 ${ONEPASSWORD_UUID}
