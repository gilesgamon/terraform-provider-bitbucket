---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_default_reviewers"
sidebar_current: "docs-bitbucket-data-repository-default-reviewers"
description: |-
  Provides information about Bitbucket repository default reviewers.
---

# bitbucket\_repository\_default\_reviewers

Provides information about Bitbucket repository default reviewers.

## Example Usage

```hcl
data "bitbucket_repository_default_reviewers" "example" {
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

* `id` - The identifier of the repository default reviewers.
* `default_reviewers` - The default reviewers. Each item contains:
    * `account_id` - The account id.
    * `account_status` - The account status.
    * `created_on` - The created on.
    * `display_name` - The display name.
    * `is_staff` - The is staff.
    * `links` - The links.
    * `nickname` - The nickname.
    * `type` - The type.
    * `username` - The username.
    * `uuid` - The uuid.
