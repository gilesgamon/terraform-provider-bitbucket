---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit"
sidebar_current: "docs-bitbucket-data-commit"
description: |-
  Provides information about Bitbucket commit.
---

# bitbucket\_commit

Provides information about Bitbucket commit.

## Example Usage

```hcl
data "bitbucket_commit" "example" {
  commit_sha = "commit_sha"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `commit_sha` - (Required) Commit SHA or branch name (e.g., 'main', 'develop', 'abc123...')
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the commit.
* `author` - The author. Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `date` - The date.
* `hash` - The hash.
* `message` - The message.
* `parents` - The parents. Each item contains:
    * `hash` - The hash.
    * `type` - The type.
