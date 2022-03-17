# You will need the vault ID and the group ID (<vault id>/<group id>).
#
# Go to the website and get the UUID from the URL or use the `op` cli:
op2 vault get test-vault
op2 group get test-group

# Import.
terraform import onepasswordorg_vault_group_access.vault0_group0 ${OP_VAULT_UUID}/${OP_GROUP_UUID}
