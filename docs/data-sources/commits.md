---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commits"
sidebar_current: "docs-bitbucket-data-commits"
description: |-
  Provides information about Bitbucket commits.
---

# bitbucket\_commits

Provides information about Bitbucket commits.

## Example Usage

```hcl
data "bitbucket_commits" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the commits.
* `commits` - The commits. Each item contains:
    * `author` - The author.
    * `date` - The date.
    * `hash` - The hash.
    * `links` - The links.
    * `message` - The message.
    * `parents` - The parents.
* `latest_commit_hash` - Hash of the most recent commit
