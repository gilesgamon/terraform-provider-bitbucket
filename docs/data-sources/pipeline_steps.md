---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_steps"
sidebar_current: "docs-bitbucket-data-pipeline-steps"
description: |-
  Provides information about Bitbucket pipeline steps.
---

# bitbucket\_pipeline\_steps

Provides information about Bitbucket pipeline steps.

## Example Usage

```hcl
data "bitbucket_pipeline_steps" "example" {
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

* `id` - The identifier of the pipeline steps.
* `steps` - The steps. Each item contains:
    * `completed_on` - The completed on.
    * `duration_in_seconds` - The duration in seconds.
    * `links` - The links.
    * `max_time` - The max time.
    * `name` - The name.
    * `script` - The script.
    * `started_on` - The started on.
    * `state` - The state.
    * `type` - The type.
    * `uuid` - The uuid.
