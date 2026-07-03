---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_schedule_executions"
sidebar_current: "docs-bitbucket-data-pipeline-schedule-executions"
description: |-
  Provides information about Bitbucket pipeline schedule executions.
---

# bitbucket\_pipeline\_schedule\_executions

Provides information about Bitbucket pipeline schedule executions.

## Example Usage

```hcl
data "bitbucket_pipeline_schedule_executions" "example" {
  repo_slug = "example-repo"
  schedule_uuid = "schedule_uuid"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) Repository slug or UUID
* `schedule_uuid` - (Required) Schedule UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline schedule executions.
* `executions` - The executions. Each item contains:
    * `created_on` - Creation timestamp
    * `pipeline` - Pipeline information
    * `uuid` - Execution UUID
