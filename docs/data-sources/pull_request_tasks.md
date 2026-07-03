---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pull_request_tasks"
sidebar_current: "docs-bitbucket-data-pull-request-tasks"
description: |-
  Provides information about Bitbucket pull request tasks.
---

# bitbucket\_pull\_request\_tasks

Provides information about Bitbucket pull request tasks.

## Example Usage

```hcl
data "bitbucket_pull_request_tasks" "example" {
  pull_request_id = "1"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pull_request_id` - (Required) Pull request ID
* `repo_slug` - (Required) Repository slug or UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pull request tasks.
* `tasks` - The tasks. Each item contains:
    * `content` - Task content
    * `created_on` - Creation timestamp
    * `creator` - Task creator
    * `id` - Task ID
    * `state` - Task state
    * `updated_on` - Last update timestamp
