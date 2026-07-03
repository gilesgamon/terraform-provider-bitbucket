---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_branch"
sidebar_current: "docs-bitbucket-data-branch"
description: |-
  Provides information about Bitbucket branch.
---

# bitbucket\_branch

Provides information about Bitbucket branch.

## Example Usage

```hcl
data "bitbucket_branch" "example" {
  branch_name = "branch_name"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `branch_name` - (Required) Branch name (e.g., 'main', 'develop', 'feature/new-feature')
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the branch.
* `name` - The name.
* `target_author` - The target author. Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `target_date` - The target date.
* `target_hash` - The target hash.
* `target_message` - The target message.
