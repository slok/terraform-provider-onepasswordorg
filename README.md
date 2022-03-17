
# terraform-provider-onepasswordorg

[![CI](https://github.com/slok/terraform-provider-onepasswordorg/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/slok/terraform-provider-onepasswordorg/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/slok/terraform-provider-onepasswordorg)](https://goreportcard.com/report/github.com/slok/terraform-provider-onepasswordorg)
[![Apache 2 licensed](https://img.shields.io/badge/license-Apache2-blue.svg)](https://raw.githubusercontent.com/slok/terraform-provider-onepasswordorg/master/LICENSE)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/slok/terraform-provider-onepasswordorg)](https://github.com/slok/terraform-provider-onepasswordorg/releases/latest)
[![Terraform regsitry](https://img.shields.io/badge/Terraform-Registry-color=green?logo=Terraform&style=flat&color=7B42BC&logoColor=white)](https://registry.terraform.io/providers/slok/onepasswordorg/latest/docs)

Terraform provider for [1password](https://1password.com) organization (e.g: Users and groups).

To manage secrets use the [official 1password](https://registry.terraform.io/providers/1Password/onepassword) provider.

## Use cases

- Create and delete users.
- Delete and create groups.
- Assign users to groups.
- Create Vaults.
- Grant fine grain vault permissions to groups.

## Requirements

- [Terraform](https://www.terraform.io/downloads.html) 1.1 or higher
- [Go](https://golang.org/doc/install) 1.18 (to build the provider plugin)
- [1password](https://app-updates.agilebits.com/product_history/CLI2) V2 CLI.

## How does it work

1password connect API doesn't allow managing resources other than secrets. Thats why this provider needs the `op` CLI.

It needs the new `op` >= `V2.x` CLI.

The `op` CLI needs a real user in onepassword to be used, so the recommendation to automate things using this provider
is to create a separate account only for automation purposes.

You will need the secret key and password of that user account.

## Terraform cloud

Terraform cloud doesn't allow installing dependencies, thats why this provider has the linux amd64 op binary embedded inside
the provider. When this provider is run from terraform cloud, it will detect, copy the op binary to "/tmp" inside terraform
cloud worker and execute that binary on the operations.

## `OP_DEVICE` error

If you are getting an error like:

```bash
cannot signin: exit status 1: [ERROR] 2022/03/14 17:13:39 No saved device ID. Set the OP_DEVICE environment variable and try again: `export OP_DEVICE=xxxxxxxxxxxxxxxxxxxxx`
```

Add the env var to your execution env with the `OP_DEVICE=xxxxxxxxxxxxxxxxxxxxx` value.

## Development

There are 2 ways while developing this provider: 

- Fake mode: Without the need for `1password` by using a fake FS storage.
- Real mode: Using a real 1password account and the `op` binary.

Both will need to build the provider. To install your plugin locally you can do `make install`, it will build and install in your `${HOME}/.terraform/plugins/...`

Note: The installation is ready for `OS_ARCH=linux_amd64`, so you make need to change the [`Makefile`](./Makefile) if using other OS.

### Fake

To enable fake storage you can use the `fake_storage_path` variable.

Example that will use `/tmp/tf-onepasswordorg-storage.json` file to store as if 1password API was called:

```terraform
provider "onepasswordorg" {
  fake_storage_path = "/tmp/tf-onepasswordorg-storage.json"
}
```

You can test this by using [fake](./examples/example) example:

```bash
make install
cd ./examples/fake
terraform init
terraform plan
```

### Real

You will need op user credentials and load them (e.g as env vars with `source ./1p-login.sh`):

```bash
export OP_ADDRESS=example.1password.com
export OP_EMAIL=bot@example.com
export OP_SECRET_KEY=XX-XXX-XXXX-XXXX-XXX
export OP_PASSWORD=xxxxxxxxxxxx
```

You can check [local](./examples/local) example.
