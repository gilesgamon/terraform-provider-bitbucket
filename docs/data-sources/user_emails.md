---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_emails"
sidebar_current: "docs-bitbucket-data-user-emails"
description: |-
  Provides information about Bitbucket user emails.
---

# bitbucket\_user\_emails

Provides information about Bitbucket user emails.

## Example Usage

```hcl
data "bitbucket_user_emails" "example" {
}
```

## Argument Reference

This data takes no arguments.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the user emails.
* `emails` - The emails. Each item contains:
    * `email` - Email address
    * `is_confirmed` - Whether this email is confirmed
    * `is_primary` - Whether this is the primary email
    * `type` - Email type
