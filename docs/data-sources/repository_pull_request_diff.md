---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pull_request_diff"
sidebar_current: "docs-bitbucket-data-repository-pull-request-diff"
description: |-
  Provides information about Bitbucket repository pull request diff.
---

# bitbucket\_repository\_pull\_request\_diff

Provides information about Bitbucket repository pull request diff.

## Example Usage

```hcl
data "bitbucket_repository_pull_request_diff" "example" {
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
* `context` - (Optional) Number of context lines to show around changes
* `path` - (Optional) Path to specific file to get diff for

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository pull request diff.
* `diff` - The diff. Each item contains:
    * `hunks` - The hunks.
    * `new` - The new.
    * `old` - The old.
