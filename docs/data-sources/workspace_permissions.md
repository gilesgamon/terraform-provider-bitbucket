---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspace_permissions"
sidebar_current: "docs-bitbucket-data-workspace-permissions"
description: |-
  Provides information about Bitbucket workspace permissions.
---

# bitbucket\_workspace\_permissions

Provides information about Bitbucket workspace permissions.

## Example Usage

```hcl
data "bitbucket_workspace_permissions" "example" {
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the workspace permissions.
* `permissions` - The permissions. Each item contains:
    * `granted_at` - The granted at.
    * `granted_by` - The granted by.
    * `group` - The group.
    * `permission` - The permission.
    * `type` - The type.
    * `user` - The user.
