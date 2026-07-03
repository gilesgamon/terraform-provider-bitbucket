---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pullrequests"
sidebar_current: "docs-bitbucket-data-pullrequests"
description: |-
  Provides information about Bitbucket pullrequests.
---

# bitbucket\_pullrequests

Provides information about Bitbucket pullrequests.

## Example Usage

```hcl
data "bitbucket_pullrequests" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.
* `author` - (Optional) Filter PRs by author username
* `destination_branch` - (Optional) Filter PRs by destination branch name
* `q` - (Optional) Search query string
* `reviewer` - (Optional) Filter PRs by reviewer username
* `sort` - (Optional) Sort field (created_on, updated_on, title, author)
* `source_branch` - (Optional) Filter PRs by source branch name
* `state` - (Optional) Filter PRs by state (open, merged, declined, superseded)

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the pullrequests.
* `pull_requests` - The pull requests. Each item contains:
    * `author` - The author.
    * `closed_by` - The closed by.
    * `closed_on` - The closed on.
    * `created_on` - The created on.
    * `description` - The description.
    * `destination` - The destination.
    * `id` - The id.
    * `merge_commit` - The merge commit.
    * `reviewers` - The reviewers.
    * `source` - The source.
    * `state` - The state.
    * `title` - The title.
    * `updated_on` - The updated on.
