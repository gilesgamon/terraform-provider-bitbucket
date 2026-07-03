---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_addon_values"
sidebar_current: "docs-bitbucket-data-repository-addon-values"
description: |-
  Provides information about Bitbucket repository addon values.
---

# bitbucket\_repository\_addon\_values

Provides information about Bitbucket repository addon values.

## Example Usage

```hcl
data "bitbucket_repository_addon_values" "example" {
  addon_key = "addon_key"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `addon_key` - (Required) The addon key.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository addon values.
* `values` - The values.
