---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_milestones"
sidebar_current: "docs-bitbucket-data-repository-milestones"
description: |-
  Provides information about Bitbucket repository milestones.
---

# bitbucket\_repository\_milestones

Provides information about Bitbucket repository milestones.

## Example Usage

```hcl
data "bitbucket_repository_milestones" "example" {
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

* `id` - The identifier of the repository milestones.
* `milestones` - The milestones. Each item contains:
    * `created_on` - The created on.
    * `description` - The description.
    * `id` - The id.
    * `links` - The links.
    * `name` - The name.
    * `release_date` - The release date.
    * `start_date` - The start date.
    * `state` - The state.
    * `updated_on` - The updated on.
