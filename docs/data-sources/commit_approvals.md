---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_approvals"
sidebar_current: "docs-bitbucket-data-commit-approvals"
description: |-
  Provides information about Bitbucket commit approvals.
---

# bitbucket\_commit\_approvals

Provides information about Bitbucket commit approvals.

## Example Usage

```hcl
data "bitbucket_commit_approvals" "example" {
  commit = "a1b2c3d4"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `commit` - (Required) The commit.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the commit approvals.
* `approvals` - The approvals. Each item contains:
    * `approved` - The approved.
    * `created_on` - The created on.
    * `links` - The links.
    * `updated_on` - The updated on.
    * `user` - The user.
    * `uuid` - The uuid.
