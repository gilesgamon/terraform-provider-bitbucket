---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspaces"
sidebar_current: "docs-bitbucket-data-workspaces"
description: |-
  Provides information about Bitbucket workspaces.
---

# bitbucket\_workspaces

Provides information about Bitbucket workspaces.

## Example Usage

```hcl
data "bitbucket_workspaces" "example" {
}
```

## Argument Reference

The following arguments are supported:

* `q` - (Optional) Search query string for workspace names

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the workspaces.
* `workspaces` - The workspaces. Each item contains:
    * `created_on` - The created on.
    * `is_private` - The is private.
    * `links` - The links.
    * `name` - The name.
    * `slug` - The slug.
    * `type` - The type.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
