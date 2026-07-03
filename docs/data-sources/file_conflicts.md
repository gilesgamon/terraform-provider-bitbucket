---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_file_conflicts"
sidebar_current: "docs-bitbucket-data-file-conflicts"
description: |-
  Provides the file conflicts for a merge base revspec in a repository
---

# bitbucket\_file\_conflicts

Retrieves the list of file conflicts that would occur when merging the given
revspec in a repository.

OAuth2 Scopes: `repository`

## Example Usage

```hcl
data "bitbucket_file_conflicts" "example" {
  workspace = "gob"
  repo_slug = "example"
  spec      = "main..feature-branch"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.
* `spec` - (Required) A merge base revspec (for example `main..feature-branch`) used to compute file conflicts.

## Attributes Reference

* `id` - The identifier of the conflict collection.
* `conflicts` - A list of file conflicts. See [File Conflict](#file-conflict) below.

### File Conflict

* `type` - The type of the conflict object.
* `path` - The path of the conflicting file.
* `scenario` - The conflict scenario (for example `content`, `rename`, `delete_modify`).
* `message` - A human-readable description of the conflict.
