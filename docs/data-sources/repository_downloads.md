---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_downloads"
sidebar_current: "docs-bitbucket-data-repository-downloads"
description: |-
  Provides information about Bitbucket repository downloads.
---

# bitbucket\_repository\_downloads

Provides information about Bitbucket repository downloads.

## Example Usage

```hcl
data "bitbucket_repository_downloads" "example" {
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

* `id` - The identifier of the repository downloads.
* `downloads` - The downloads. Each item contains:
    * `created_on` - The created on.
    * `downloads` - The downloads.
    * `links` - The links.
    * `name` - The name.
    * `size` - The size.
