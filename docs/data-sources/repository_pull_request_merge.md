---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pull_request_merge"
sidebar_current: "docs-bitbucket-data-repository-pull-request-merge"
description: |-
  Provides information about Bitbucket repository pull request merge.
---

# bitbucket\_repository\_pull\_request\_merge

Provides information about Bitbucket repository pull request merge.

## Example Usage

```hcl
data "bitbucket_repository_pull_request_merge" "example" {
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

* `id` - The identifier of the repository pull request merge.
* `close_source_branch` - The close source branch.
* `destination` - The destination.
* `merge_commit` - The merge commit.
* `merge_status` - The merge status.
* `merge_strategy` - The merge strategy.
* `source` - The source.
