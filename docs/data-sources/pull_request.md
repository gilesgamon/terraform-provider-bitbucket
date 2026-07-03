---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pull_request"
sidebar_current: "docs-bitbucket-data-pull-request"
description: |-
  Provides information about Bitbucket pull request.
---

# bitbucket\_pull\_request

Provides information about Bitbucket pull request.

## Example Usage

```hcl
data "bitbucket_pull_request" "example" {
  pull_request_id = "1"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pull_request_id` - (Required) Pull request ID (number)
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pull request.
* `author` - The author. Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `created_date` - The created date.
* `description` - The description.
* `destination` - The destination. Each item contains:
    * `branch` - The branch.
    * `commit` - The commit.
    * `repository` - The repository.
* `merge_commit` - The merge commit.
* `source` - The source. Each item contains:
    * `branch` - The branch.
    * `commit` - The commit.
    * `repository` - The repository.
* `state` - The state.
* `title` - The title.
* `updated_date` - The updated date.
