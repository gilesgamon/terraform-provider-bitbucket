---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_files"
sidebar_current: "docs-bitbucket-data-repository-files"
description: |-
  Provides information about Bitbucket repository files.
---

# bitbucket\_repository\_files

Provides information about Bitbucket repository files.

## Example Usage

```hcl
data "bitbucket_repository_files" "example" {
  ref = "ref"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `ref` - (Required) Branch, tag, or commit hash
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.
* `path` - (Optional) Path to directory or file

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository files.
* `files` - The files. Each item contains:
    * `hash` - The hash.
    * `links` - The links.
    * `path` - The path.
    * `size` - The size.
    * `type` - The type.
