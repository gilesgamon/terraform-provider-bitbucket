---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_workspace_repository_permissions"
sidebar_current: "docs-bitbucket-data-user-workspace-repository-permissions"
description: |-
  Provides the current user's repository permissions within a workspace
---

# bitbucket\_user\_workspace\_repository\_permissions

Retrieves the currently authenticated user's repository permissions within a
given workspace.

OAuth2 Scopes: `repository`

## Example Usage

```hcl
data "bitbucket_user_workspace_repository_permissions" "example" {
  workspace = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `q` - (Optional) Query string to narrow down the response.
* `sort` - (Optional) Field by which the results should be sorted.

## Attributes Reference

* `id` - The identifier of the permissions collection.
* `permissions` - A list of repository permissions. See [Permission](#permission) below.

### Permission

* `permission` - The permission level (`read`, `write`, `admin` or `none`).
* `repository_uuid` - The UUID of the repository.
* `repository_name` - The name of the repository.
* `repository_full_name` - The full name of the repository (`workspace/repo`).
