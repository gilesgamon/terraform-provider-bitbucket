---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_branch_merge_base"
sidebar_current: "docs-bitbucket-data-branch-merge-base"
description: |-
  Provides the common ancestor commit between two revisions in a Bitbucket repository
---

# bitbucket\_branch\_merge\_base

Retrieves the common ancestor (merge base) commit between two revisions —
branches, tags or commit hashes — in a repository.

OAuth2 Scopes: `repository`

## Example Usage

```hcl
data "bitbucket_branch_merge_base" "example" {
  workspace = "example-workspace"
  repo_slug = "example-repo"
  source    = "main"
  target    = "feature-branch"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.
* `source` - (Required) The first revision (branch name, tag or commit hash).
* `target` - (Required) The second revision (branch name, tag or commit hash).

## Attributes Reference

The common ancestor is returned as a commit:

* `id` - The identifier of the merge base lookup.
* `hash` - The hash of the common ancestor commit.
* `message` - The commit message.
* `date` - The commit date.
* `author` - The commit author. Each item contains `username`, `display_name` and `uuid`.
* `parents` - The parent commits. Each item contains `hash` and `type`.
