---
# generated by https://github.com/hashicorp/terraform-plugin-docs
page_title: "onepasswordorg_user Data Source - terraform-provider-onepasswordorg"
subcategory: ""
description: |-
  Provides information about a 1password user.
---

# onepasswordorg_user (Data Source)

Provides information about a 1password user.

## Example Usage

```terraform
data "onepasswordorg_user" "test" {
  email = "user0@slok.dev"
}

output "user_test" {
  value = data.onepasswordorg_user.test
}
```

<!-- schema generated by tfplugindocs -->
## Schema

### Required

- **email** (String) The email of the user.

### Read-Only

- **id** (String) The ID of this resource.
- **name** (String)

