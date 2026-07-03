---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_issue_comments"
sidebar_current: "docs-bitbucket-data-issue-comments"
description: |-
  Provides information about Bitbucket issue comments.
---

# bitbucket\_issue\_comments

Provides information about Bitbucket issue comments.

## Example Usage

```hcl
data "bitbucket_issue_comments" "example" {
  issue_id = "issue_id"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `issue_id` - (Required) The issue id.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the issue comments.
* `comments` - The comments. Each item contains:
    * `content` - The content.
    * `created_on` - The created on.
    * `id` - The id.
    * `links` - The links.
    * `updated_on` - The updated on.
    * `user` - The user.
