---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_permissions"
sidebar_current: "docs-bitbucket-data-repository-permissions"
description: |-
  Provides information about Bitbucket repository permissions.
---

# bitbucket\_repository\_permissions

Provides information about Bitbucket repository permissions.

## Example Usage

```hcl
data "bitbucket_repository_permissions" "example" {
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

* `id` - The identifier of the repository permissions.
* `permissions` - The permissions. Each item contains:
    * `granted_at` - The granted at.
    * `granted_by` - The granted by.
    * `permission` - The permission.
    * `repository` - The repository.
    * `type` - The type.
    * `user` - The user.
