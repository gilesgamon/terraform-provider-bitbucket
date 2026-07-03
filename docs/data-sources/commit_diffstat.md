---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_diffstat"
sidebar_current: "docs-bitbucket-data-commit-diffstat"
description: |-
  Provides information about Bitbucket commit diffstat.
---

# bitbucket\_commit\_diffstat

Provides information about Bitbucket commit diffstat.

## Example Usage

```hcl
data "bitbucket_commit_diffstat" "example" {
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

* `id` - The identifier of the commit diffstat.
* `diffstat` - The diffstat. Each item contains:
    * `deleted_file` - The deleted file.
    * `lines_added` - The lines added.
    * `lines_removed` - The lines removed.
    * `new_file` - The new file.
    * `new_path` - The new path.
    * `old_path` - The old path.
    * `renamed_file` - The renamed file.
    * `status` - The status.
    * `type` - The type.
