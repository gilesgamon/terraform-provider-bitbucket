---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_issue_changes"
sidebar_current: "docs-bitbucket-data-repository-issue-changes"
description: |-
  Provides information about Bitbucket repository issue changes.
---

# bitbucket\_repository\_issue\_changes

Provides information about Bitbucket repository issue changes.

## Example Usage

```hcl
data "bitbucket_repository_issue_changes" "example" {
  issue_id = "issue_id"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `issue_id` - (Required) Issue ID
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository issue changes.
* `changes` - The changes. Each item contains:
    * `changes` - The changes.
    * `created_on` - The created on.
    * `id` - The id.
    * `links` - The links.
    * `type` - The type.
    * `user` - The user.
