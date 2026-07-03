---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pipelines"
sidebar_current: "docs-bitbucket-data-pipelines"
description: |-
  Provides information about Bitbucket pipelines.
---

# bitbucket\_pipelines

Provides information about Bitbucket pipelines.

## Example Usage

```hcl
data "bitbucket_pipelines" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.
* `page` - (Optional) Page number for pagination
* `state` - (Optional) Filter pipelines by state (pending, in_progress, completed, error, stopped)
* `target` - (Optional) Filter pipelines by target (commit, tag, branch, custom)

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pipelines.
* `pipelines` - The pipelines. Each item contains:
    * `build_number` - The build number.
    * `build_seconds_used` - The build seconds used.
    * `completed_on` - The completed on.
    * `created_on` - The created on.
    * `duration_in_seconds` - The duration in seconds.
    * `expired` - The expired.
    * `first_successful` - The first successful.
    * `repository` - The repository.
    * `state` - The state.
    * `target` - The target.
    * `trigger` - The trigger.
    * `uuid` - The uuid.
