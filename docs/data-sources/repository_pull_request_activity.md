---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pull_request_activity"
sidebar_current: "docs-bitbucket-data-repository-pull-request-activity"
description: |-
  Provides information about Bitbucket repository pull request activity.
---

# bitbucket\_repository\_pull\_request\_activity

Provides information about Bitbucket repository pull request activity.

## Example Usage

```hcl
data "bitbucket_repository_pull_request_activity" "example" {
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

* `id` - The identifier of the repository pull request activity.
* `activity` - The activity. Each item contains:
    * `approved` - The approved.
    * `changes_requested` - The changes requested.
    * `comment` - The comment.
    * `created_on` - The created on.
    * `id` - The id.
    * `links` - The links.
    * `type` - The type.
    * `update` - The update.
    * `user` - The user.
