---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_pullrequests"
sidebar_current: "docs-bitbucket-data-commit-pullrequests"
description: |-
  Provides information about Bitbucket commit pullrequests.
---

# bitbucket\_commit\_pullrequests

Provides information about Bitbucket commit pullrequests.

## Example Usage

```hcl
data "bitbucket_commit_pullrequests" "example" {
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

* `id` - The identifier of the commit pullrequests.
* `pullrequests` - The pullrequests. Each item contains:
    * `author` - The author.
    * `created_on` - The created on.
    * `description` - The description.
    * `destination` - The destination.
    * `id` - The id.
    * `links` - The links.
    * `source` - The source.
    * `state` - The state.
    * `title` - The title.
    * `updated_on` - The updated on.
