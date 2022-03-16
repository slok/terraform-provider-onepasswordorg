# Changelog

## [Unreleased]

### Added

- Vault resource.
- Vault data source.

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

[unreleased]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.2.0...HEAD
[v0.2.0]: https://github.com/slok/terraform-provider-onepasswordorg/compare/v0.1.0...v0.2.0
[v0.1.0]: https://github.com/slok/terraform-provider-onepasswordorg/releases/tag/v0.1.0
