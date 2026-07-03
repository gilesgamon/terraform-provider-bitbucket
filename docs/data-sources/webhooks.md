---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_webhooks"
sidebar_current: "docs-bitbucket-data-webhooks"
description: |-
  Provides information about Bitbucket webhooks.
---

# bitbucket\_webhooks

Provides information about Bitbucket webhooks.

## Example Usage

```hcl
data "bitbucket_webhooks" "example" {
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the webhooks.
* `webhooks` - The webhooks. Each item contains:
    * `active` - The active.
    * `created_on` - The created on.
    * `description` - The description.
    * `events` - The events.
    * `links` - The links.
    * `skip_cert_verification` - The skip cert verification.
    * `subject` - The subject.
    * `updated_on` - The updated on.
    * `url` - The url.
    * `uuid` - The uuid.
