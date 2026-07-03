---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_project_permissions"
sidebar_current: "docs-bitbucket-data-project-permissions"
description: |-
  Provides information about Bitbucket project permissions.
---

# bitbucket\_project\_permissions

Provides information about Bitbucket project permissions.

## Example Usage

```hcl
data "bitbucket_project_permissions" "example" {
  project_key = "PROJ"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `project_key` - (Required) The project key.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the project permissions.
* `permissions` - The permissions. Each item contains:
    * `granted_at` - The granted at.
    * `granted_by` - The granted by.
    * `group` - The group.
    * `permission` - The permission.
    * `type` - The type.
    * `user` - The user.
