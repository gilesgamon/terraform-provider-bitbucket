---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pull_request_merge_task_status"
sidebar_current: "docs-bitbucket-data-pull-request-merge-task-status"
description: |-
  Provides information about Bitbucket pull request merge task status.
---

# bitbucket\_pull\_request\_merge\_task\_status

Provides information about Bitbucket pull request merge task status.

## Example Usage

```hcl
data "bitbucket_pull_request_merge_task_status" "example" {
  pull_request_id = "1"
  repo_slug = "example-repo"
  task_id = "task_id"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pull_request_id` - (Required) Pull request ID
* `repo_slug` - (Required) Repository slug or UUID
* `task_id` - (Required) Task ID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pull request merge task status.
* `created_on` - Creation timestamp
* `status` - Merge task status
* `updated_on` - Last update timestamp
