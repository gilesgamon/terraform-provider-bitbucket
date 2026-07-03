---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_file_history"
sidebar_current: "docs-bitbucket-data-repository-file-history"
description: |-
  Provides information about Bitbucket repository file history.
---

# bitbucket\_repository\_file\_history

Provides information about Bitbucket repository file history.

## Example Usage

```hcl
data "bitbucket_repository_file_history" "example" {
  path = "path"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `path` - (Required) Path to the file
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.
* `revision` - (Optional) Revision (commit hash) to get history for

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository file history.
* `history` - The history. Each item contains:
    * `commit` - The commit.
    * `links` - The links.
    * `path` - The path.
    * `size` - The size.
    * `type` - The type.
