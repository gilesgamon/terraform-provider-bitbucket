---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_test_case_reasons"
sidebar_current: "docs-bitbucket-data-pipeline-test-case-reasons"
description: |-
  Provides information about Bitbucket pipeline test case reasons.
---

# bitbucket\_pipeline\_test\_case\_reasons

Provides information about Bitbucket pipeline test case reasons.

## Example Usage

```hcl
data "bitbucket_pipeline_test_case_reasons" "example" {
  pipeline_uuid = "pipeline_uuid"
  repo_slug = "example-repo"
  step_uuid = "step_uuid"
  test_case_uuid = "test_case_uuid"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pipeline_uuid` - (Required) Pipeline UUID
* `repo_slug` - (Required) Repository slug or UUID
* `step_uuid` - (Required) Step UUID
* `test_case_uuid` - (Required) Test case UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipeline test case reasons.
* `reasons` - The reasons. Each item contains:
    * `created_on` - Creation timestamp
    * `description` - Reason description
    * `name` - Reason name
    * `type` - Reason type
