data "onepasswordorg_vault" "test" {
  name = "test-vault"
}

output "vault_test" {
  value = data.onepasswordorg_vault.test
}
