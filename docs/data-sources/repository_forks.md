---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_forks"
sidebar_current: "docs-bitbucket-data-repository-forks"
description: |-
  Provides information about Bitbucket repository forks.
---

# bitbucket\_repository\_forks

Provides information about Bitbucket repository forks.

## Example Usage

```hcl
data "bitbucket_repository_forks" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.
* `q` - (Optional) Search query string for repository names

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository forks.
* `forks` - The forks. Each item contains:
    * `created_on` - The created on.
    * `description` - The description.
    * `fork_policy` - The fork policy.
    * `full_name` - The full name.
    * `is_private` - The is private.
    * `language` - The language.
    * `mainbranch` - The mainbranch.
    * `name` - The name.
    * `project` - The project.
    * `size` - The size.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
    * `workspace` - The workspace.
