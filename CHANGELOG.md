# Changelog

## [Unreleased]

## [v0.6.0] - 2024-10-22

### Changed

- Update Go version to v1.23.
- Update terraform related dependencies.
- Remove all `terraform-plugin-sdk` usage in favor of the new `terraform-plugin-framework` libs.
- Updated TFE embedded `op` CLI to 2.30.0

## [v0.5.0] - 2022-07-30

### Changed

- Updated terraform dependencies (SDKs, plugins...).
- Updated TFE embedded `op` CLI to 2.6.0.

## [v0.4.0] - 2022-03-18

### Added

- Vault user access with fine grain permissions.

## [v0.3.0] - 2022-03-17

### Added

- Vault resource.
- Vault data source.
- Provider option to allow selecting a specific op cli in a specific path.
- Vault group access with fine grain permissions.

## Changed

- `onepasswordorg_group.description` now is optional and defaults to `Managed by Terraform`.
- `onepasswordorg_group_member.role` now is optional and defaults to `member`.
- `onepasswordorg_group_member.role` validates that is not set as empty string.
- `onepasswordorg_group_member.user_id` validates that is not set as empty string.
- `onepasswordorg_group_member.group_id` validates that is not set as empty string.
- `onepasswordorg_group.name` validates that is not set as empty string.
- `onepasswordorg_user.name` validates that is not set as empty string.
- `onepasswordorg_email.name` validates that is not set as empty string.
- Use go 1.18.
- Use v2.0.0 op cli.

## Deleted

- Support v1.x op cli

## [v0.2.0] - 2022-03-15

### Added

- User data source.
- Group data source.

## Changed

- Group member now doesn't need two tf applies to be have a role different than "member"

## [v0.1.0] - 2022-03-15

### Added

- User resource.
- Group resource.
- Group member resource.
- Fake 1password storage.
- Real 1password storage.
- Terraform registry release.

[unreleased]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.6.0...HEAD
[v0.5.0]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.6.0...v0.6.0
[v0.5.0]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.4.0...v0.5.0
[v0.4.0]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.3.0...v0.4.0
[v0.3.0]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.2.0...v0.3.0
[v0.2.0]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/slok/terraform-provider-onepasswordorg/releases/tag/v0.1.0
