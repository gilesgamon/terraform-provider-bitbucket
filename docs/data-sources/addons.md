---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_addons"
sidebar_current: "docs-bitbucket-data-addons"
description: |-
  Provides information about Bitbucket addons.
---

# bitbucket\_addons

Provides information about Bitbucket addons.

## Example Usage

```hcl
data "bitbucket_addons" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the addons.
* `addons` - The addons. Each item contains:
    * `addon_key` - The addon key.
    * `app_info` - The app info.
    * `description` - The description.
    * `enabled` - The enabled.
    * `installed` - The installed.
    * `links` - The links.
    * `name` - The name.
    * `vendor` - The vendor.
