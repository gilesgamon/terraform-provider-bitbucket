---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_logs"
sidebar_current: "docs-bitbucket-data-pipeline-logs"
description: |-
  Provides information about Bitbucket pipeline logs.
---

# bitbucket\_pipeline\_logs

Provides information about Bitbucket pipeline logs.

## Example Usage

```hcl
data "bitbucket_pipeline_logs" "example" {
  pipeline_uuid = "pipeline_uuid"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_uuid` - (Required) The pipeline uuid.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline logs.
* `logs` - The logs. Each item contains:
    * `created_on` - The created on.
    * `level` - The level.
    * `message` - The message.
    * `step` - The step.
    * `updated_on` - The updated on.
