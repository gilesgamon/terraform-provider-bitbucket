---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_statuses"
sidebar_current: "docs-bitbucket-data-commit-statuses"
description: |-
  Provides information about Bitbucket commit statuses.
---

# bitbucket\_commit\_statuses

Provides information about Bitbucket commit statuses.

## Example Usage

```hcl
data "bitbucket_commit_statuses" "example" {
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

* `id` - The identifier of the commit statuses.
* `statuses` - The statuses. Each item contains:
    * `created_on` - The created on.
    * `description` - The description.
    * `key` - The key.
    * `links` - The links.
    * `name` - The name.
    * `refname` - The refname.
    * `state` - The state.
    * `updated_on` - The updated on.
    * `url` - The url.
    * `uuid` - The uuid.
