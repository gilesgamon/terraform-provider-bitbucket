---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline"
sidebar_current: "docs-bitbucket-data-pipeline"
description: |-
  Provides information about Bitbucket pipeline.
---

# bitbucket\_pipeline

Provides information about Bitbucket pipeline.

## Example Usage

```hcl
data "bitbucket_pipeline" "example" {
  pipeline_number = "pipeline_number"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_number` - (Required) Pipeline number to retrieve
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline.
* `build_number` - The build number.
* `completed_on` - The completed on.
* `created_on` - The created on.
* `state` - The state.
* `target` - The target. Each item contains:
    * `hash` - The hash.
    * `ref_name` - The ref name.
    * `type` - The type.
* `trigger` - The trigger. Each item contains:
    * `type` - The type.
    * `user` - The user.
