---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_components"
sidebar_current: "docs-bitbucket-data-repository-components"
description: |-
  Provides information about Bitbucket repository components.
---

# bitbucket\_repository\_components

Provides information about Bitbucket repository components.

## Example Usage

```hcl
data "bitbucket_repository_components" "example" {
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

* `id` - The identifier of the repository components.
* `components` - The components. Each item contains:
    * `assignee` - The assignee.
    * `created_on` - The created on.
    * `description` - The description.
    * `id` - The id.
    * `links` - The links.
    * `name` - The name.
    * `updated_on` - The updated on.
