---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_schedules"
sidebar_current: "docs-bitbucket-data-repository-pipeline-schedules"
description: |-
  Provides information about Bitbucket repository pipeline schedules.
---

# bitbucket\_repository\_pipeline\_schedules

Provides information about Bitbucket repository pipeline schedules.

## Example Usage

```hcl
data "bitbucket_repository_pipeline_schedules" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository pipeline schedules.
* `schedules` - The schedules. Each item contains:
    * `created_on` - The created on.
    * `cron_pattern` - The cron pattern.
    * `enabled` - The enabled.
    * `links` - The links.
    * `name` - The name.
    * `next_run` - The next run.
    * `target` - The target.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
