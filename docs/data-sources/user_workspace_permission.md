---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_workspace_permission"
sidebar_current: "docs-bitbucket-data-user-workspace-permission"
description: |-
  Provides the current user's membership for a workspace
---

# bitbucket\_user\_workspace\_permission

Retrieves the currently authenticated user's membership for a given workspace.

OAuth2 Scopes: `account`

## Example Usage

```hcl
data "bitbucket_user_workspace_permission" "example" {
  workspace = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.

## Attributes Reference

* `id` - The identifier of the membership.
* `user_uuid` - The UUID of the user.
* `user_nickname` - The nickname of the user.
* `workspace_uuid` - The UUID of the workspace.
* `workspace_slug` - The slug of the workspace.
