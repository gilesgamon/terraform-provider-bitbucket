---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_workspaces"
sidebar_current: "docs-bitbucket-data-user-workspaces"
description: |-
  Provides the workspaces the current user is a member of
---

# bitbucket\_user\_workspaces

Retrieves the list of workspaces that the currently authenticated user is a
member of, together with the user's administrator permission on each.

OAuth2 Scopes: `account`

## Example Usage

```hcl
data "bitbucket_user_workspaces" "example" {
  administrator = true
}
```

## Argument Reference

The following arguments are supported:

* `sort` - (Optional) Field by which the results should be sorted.
* `administrator` - (Optional) Only return workspaces where the current user is an administrator.

## Attributes Reference

* `id` - The identifier of the workspaces collection.
* `workspaces` - A list of workspaces. See [Workspace](#workspace) below.

### Workspace

* `administrator` - Whether the current user is an administrator of the workspace.
* `uuid` - The workspace's immutable id.
* `slug` - The short label that identifies the workspace.
* `name` - The workspace name.
* `type` - The type of the object.
