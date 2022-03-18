# You will need the vault ID and the user ID (<vault id>/<user id>).
#
# Go to the website and get the UUID from the URL or use the `op` cli:
op2 vault get test-vault
op2 user get test-user

# Import.
terraform import onepasswordorg_vault_user_access.vault0_user0 ${OP_VAULT_UUID}/${OP_USER_UUID}
