---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_diff"
sidebar_current: "docs-bitbucket-data-commit-diff"
description: |-
  Provides information about Bitbucket commit diff.
---

# bitbucket\_commit\_diff

Provides information about Bitbucket commit diff.

## Example Usage

```hcl
data "bitbucket_commit_diff" "example" {
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

* `id` - The identifier of the commit diff.
* `diff` - The diff. Each item contains:
    * `deleted_file` - The deleted file.
    * `hunks` - The hunks.
    * `lines_added` - The lines added.
    * `lines_removed` - The lines removed.
    * `new_file` - The new file.
    * `new_path` - The new path.
    * `old_path` - The old path.
    * `renamed_file` - The renamed file.
    * `similarity` - The similarity.
    * `status` - The status.
