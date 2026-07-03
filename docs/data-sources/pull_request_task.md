---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pull_request_task"
sidebar_current: "docs-bitbucket-data-pull-request-task"
description: |-
  Provides information about Bitbucket pull request task.
---

# bitbucket\_pull\_request\_task

Provides information about Bitbucket pull request task.

## Example Usage

```hcl
data "bitbucket_pull_request_task" "example" {
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

* `id` - The identifier of the pull request task.
* `content` - Task content
* `created_on` - Creation timestamp
* `creator` - Task creator Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `state` - Task state
* `updated_on` - Last update timestamp
