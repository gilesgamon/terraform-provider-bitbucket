---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_comments"
sidebar_current: "docs-bitbucket-data-commit-comments"
description: |-
  Provides information about Bitbucket commit comments.
---

# bitbucket\_commit\_comments

Provides information about Bitbucket commit comments.

## Example Usage

```hcl
data "bitbucket_commit_comments" "example" {
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

* `id` - The identifier of the commit comments.
* `comments` - The comments. Each item contains:
    * `content` - The content.
    * `created_on` - The created on.
    * `id` - The id.
    * `links` - The links.
    * `updated_on` - The updated on.
    * `user` - The user.
