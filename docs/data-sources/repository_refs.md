---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_refs"
sidebar_current: "docs-bitbucket-data-repository-refs"
description: |-
  Provides information about Bitbucket repository refs.
---

# bitbucket\_repository\_refs

Provides information about Bitbucket repository refs.

## Example Usage

```hcl
data "bitbucket_repository_refs" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.
* `q` - (Optional) Search query string for ref names
* `sort` - (Optional) Sort order (name, -name, target, -target)

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository refs.
* `refs` - The refs. Each item contains:
    * `links` - The links.
    * `name` - The name.
    * `target` - The target.
    * `type` - The type.
