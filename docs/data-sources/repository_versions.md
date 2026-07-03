---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_versions"
sidebar_current: "docs-bitbucket-data-repository-versions"
description: |-
  Provides information about Bitbucket repository versions.
---

# bitbucket\_repository\_versions

Provides information about Bitbucket repository versions.

## Example Usage

```hcl
data "bitbucket_repository_versions" "example" {
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

* `id` - The identifier of the repository versions.
* `versions` - The versions. Each item contains:
    * `created_on` - The created on.
    * `description` - The description.
    * `id` - The id.
    * `links` - The links.
    * `name` - The name.
    * `released` - The released.
    * `released_on` - The released on.
    * `updated_on` - The updated on.
