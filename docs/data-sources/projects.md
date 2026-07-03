---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_projects"
sidebar_current: "docs-bitbucket-data-projects"
description: |-
  Provides information about Bitbucket projects.
---

# bitbucket\_projects

Provides information about Bitbucket projects.

## Example Usage

```hcl
data "bitbucket_projects" "example" {
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace.
* `q` - (Optional) Search query string for project names

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the projects.
* `projects` - The projects. Each item contains:
    * `created_on` - The created on.
    * `description` - The description.
    * `is_private` - The is private.
    * `key` - The key.
    * `links` - The links.
    * `name` - The name.
    * `owner` - The owner.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
