---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_issue_import"
sidebar_current: "docs-bitbucket-data-repository-issue-import"
description: |-
  Provides information about Bitbucket repository issue import.
---

# bitbucket\_repository\_issue\_import

Provides information about Bitbucket repository issue import.

## Example Usage

```hcl
data "bitbucket_repository_issue_import" "example" {
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

* `id` - The identifier of the repository issue import.
* `created_on` - Creation timestamp
* `import_status` - Import status
* `import_url` - Import URL
* `updated_on` - Last update timestamp
