# terraform-provider-onepasswordorg

Terraform provider for 1password organization (Users and groups).

To manage secrets use the [official 1password](https://registry.terraform.io/providers/1Password/onepassword) provider.

## Use cases

- Create and delete users.
- Delete and create groups.
- Assign users to groups.

## How does it work

1password connect API doesn't allow managing resources other than secrets. Thats we this provider needs the `op` CLI.

The `op` CLI needs a real user in onepassword to be used, so the recomendation to automate things using this provider
is to create a separate account only for automation purposes.

You will need the secret key and password of that user account.

## Terraform cloud

Terraform cloud doesn't allow installing dependencies, thats why this provider has the linux amd64 op binary embedded inside
the provider. When this provider is run from terraform cloud, it will detect, copy the op binary to "/tmp" inside terraform
cloud worker and execute that binary on the operations.
