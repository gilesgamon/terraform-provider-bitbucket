---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pull_request_approvals"
sidebar_current: "docs-bitbucket-data-repository-pull-request-approvals"
description: |-
  Provides information about Bitbucket repository pull request approvals.
---

# bitbucket\_repository\_pull\_request\_approvals

Provides information about Bitbucket repository pull request approvals.

## Example Usage

```hcl
data "bitbucket_repository_pull_request_approvals" "example" {
  pull_request_id = "1"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `pull_request_id` - (Required) Pull request ID
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository pull request approvals.
* `approvals` - The approvals. Each item contains:
    * `approved_on` - The approved on.
    * `links` - The links.
    * `role` - The role.
    * `user` - The user.
