---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_test_cases"
sidebar_current: "docs-bitbucket-data-pipeline-test-cases"
description: |-
  Provides information about Bitbucket pipeline test cases.
---

# bitbucket\_pipeline\_test\_cases

Provides information about Bitbucket pipeline test cases.

## Example Usage

```hcl
data "bitbucket_pipeline_test_cases" "example" {
  pipeline_uuid = "pipeline_uuid"
  repo_slug = "example-repo"
  step_uuid = "step_uuid"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_uuid` - (Required) Pipeline UUID
* `repo_slug` - (Required) Repository slug or UUID
* `step_uuid` - (Required) Step UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline test cases.
* `test_cases` - The test cases. Each item contains:
    * `classname` - Test case class name
    * `created_on` - Creation timestamp
    * `duration` - Test case duration in seconds
    * `file` - Test case file
    * `name` - Test case name
    * `result` - Test case result
    * `uuid` - Test case UUID
