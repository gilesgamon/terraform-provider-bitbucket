---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_issue_votes"
sidebar_current: "docs-bitbucket-data-repository-issue-votes"
description: |-
  Provides information about Bitbucket repository issue votes.
---

# bitbucket\_repository\_issue\_votes

Provides information about Bitbucket repository issue votes.

## Example Usage

```hcl
data "bitbucket_repository_issue_votes" "example" {
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

* `id` - The identifier of the repository issue votes.
* `votes` - The votes. Each item contains:
    * `created_on` - The created on.
    * `links` - The links.
    * `user` - The user.
