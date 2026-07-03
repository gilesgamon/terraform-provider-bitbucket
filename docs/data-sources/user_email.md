---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_email"
sidebar_current: "docs-bitbucket-data-user-email"
description: |-
  Provides information about Bitbucket user email.
---

# bitbucket\_user\_email

Provides information about Bitbucket user email.

## Example Usage

```hcl
data "bitbucket_user_email" "example" {
  email = "email"
}
```

## Argument Reference

The following arguments are supported:

* `email` - (Required) Email address

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the user email.
* `is_confirmed` - Whether this email is confirmed
* `is_primary` - Whether this is the primary email
* `type` - Email type
