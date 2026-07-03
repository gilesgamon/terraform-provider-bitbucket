---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_stop"
sidebar_current: "docs-bitbucket-resource-pipeline-stop"
description: |-
  Provides a Bitbucket pipeline stop resource.
---

# bitbucket\_pipeline\_stop

Provides a Bitbucket pipeline stop resource.

## Example Usage

```hcl
resource "bitbucket_pipeline_stop" "example" {
  pipeline_uuid = "pipeline_uuid"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_uuid` - (Required) Pipeline UUID
* `repo_slug` - (Required) Repository slug or UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline stop.
* `stopped` - Whether the pipeline was successfully stopped
