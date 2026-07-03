---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_reports"
sidebar_current: "docs-bitbucket-data-commit-reports"
description: |-
  Provides information about Bitbucket commit reports.
---

# bitbucket\_commit\_reports

Provides information about Bitbucket commit reports.

## Example Usage

```hcl
data "bitbucket_commit_reports" "example" {
  commit = "a1b2c3d4"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `commit` - (Required) The commit.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the commit reports.
* `reports` - The reports. Each item contains:
    * `created_on` - The created on.
    * `data` - The data.
    * `details` - The details.
    * `external_id` - The external id.
    * `link` - The link.
    * `reporter` - The reporter.
    * `result` - The result.
    * `severity` - The severity.
    * `title` - The title.
    * `type` - The type.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
