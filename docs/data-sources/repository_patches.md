---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_patches"
sidebar_current: "docs-bitbucket-data-repository-patches"
description: |-
  Provides information about Bitbucket repository patches.
---

# bitbucket\_repository\_patches

Provides information about Bitbucket repository patches.

## Example Usage

```hcl
data "bitbucket_repository_patches" "example" {
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

* `id` - The identifier of the repository patches.
* `patches` - The patches. Each item contains:
    * `created_on` - The created on.
    * `links` - The links.
    * `name` - The name.
    * `size` - The size.
