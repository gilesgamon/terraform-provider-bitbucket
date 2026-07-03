---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipeline_test_reports"
sidebar_current: "docs-bitbucket-data-pipeline-test-reports"
description: |-
  Provides information about Bitbucket pipeline test reports.
---

# bitbucket\_pipeline\_test\_reports

Provides information about Bitbucket pipeline test reports.

## Example Usage

```hcl
data "bitbucket_pipeline_test_reports" "example" {
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

* `id` - The identifier of the pipeline test reports.
* `test_reports` - The test reports. Each item contains:
    * `duration` - The duration.
    * `failed_tests` - The failed tests.
    * `name` - The name.
    * `passed_tests` - The passed tests.
    * `skipped_tests` - The skipped tests.
    * `total_tests` - The total tests.
