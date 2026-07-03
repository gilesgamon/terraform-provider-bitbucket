---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_users"
sidebar_current: "docs-bitbucket-data-users"
description: |-
  Provides information about Bitbucket users.
---

# bitbucket\_users

Provides information about Bitbucket users.

## Example Usage

```hcl
data "bitbucket_users" "example" {
}
```

## Argument Reference

The following arguments are supported:

* `q` - (Optional) Search query string for usernames or display names

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the users.
* `users` - The users. Each item contains:
    * `account_id` - The account id.
    * `account_status` - The account status.
    * `created_on` - The created on.
    * `display_name` - The display name.
    * `is_staff` - The is staff.
    * `links` - The links.
    * `nickname` - The nickname.
    * `type` - The type.
    * `username` - The username.
    * `uuid` - The uuid.
