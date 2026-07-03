---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pull_request_comments"
sidebar_current: "docs-bitbucket-data-repository-pull-request-comments"
description: |-
  Provides information about Bitbucket repository pull request comments.
---

# bitbucket\_repository\_pull\_request\_comments

Provides information about Bitbucket repository pull request comments.

## Example Usage

```hcl
data "bitbucket_repository_pull_request_comments" "example" {
  pull_request_id = "1"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pull_request_id` - (Required) Pull request ID
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository pull request comments.
* `comments` - The comments. Each item contains:
    * `content` - The content.
    * `created_on` - The created on.
    * `deleted` - The deleted.
    * `id` - The id.
    * `inline` - The inline.
    * `links` - The links.
    * `parent` - The parent.
    * `updated_on` - The updated on.
    * `user` - The user.
