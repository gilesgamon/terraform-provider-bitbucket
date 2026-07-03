---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_issue_export"
sidebar_current: "docs-bitbucket-data-repository-issue-export"
description: |-
  Provides information about Bitbucket repository issue export.
---

# bitbucket\_repository\_issue\_export

Provides information about Bitbucket repository issue export.

## Example Usage

```hcl
data "bitbucket_repository_issue_export" "example" {
  export_id = "export_id"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `export_id` - (Required) The export id.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository issue export.
* `export` - The export. Each item contains:
    * `created_on` - The created on.
    * `download_url` - The download url.
    * `links` - The links.
    * `status` - The status.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
