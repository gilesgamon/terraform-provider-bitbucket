---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_variables"
sidebar_current: "docs-bitbucket-data-repository-variables"
description: |-
  Provides information about Bitbucket repository variables.
---

# bitbucket\_repository\_variables

Provides information about Bitbucket repository variables.

## Example Usage

```hcl
data "bitbucket_repository_variables" "example" {
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

* `id` - The identifier of the repository variables.
* `variables` - The variables. Each item contains:
    * `created_on` - The created on.
    * `key` - The key.
    * `secured` - The secured.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
    * `value` - The value.
