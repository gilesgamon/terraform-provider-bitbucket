---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_build_number"
sidebar_current: "docs-bitbucket-data-pipeline-build-number"
description: |-
  Provides information about Bitbucket pipeline build number.
---

# bitbucket\_pipeline\_build\_number

Provides information about Bitbucket pipeline build number.

## Example Usage

```hcl
data "bitbucket_pipeline_build_number" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) Repository slug or UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline build number.
* `next` - Next build number
