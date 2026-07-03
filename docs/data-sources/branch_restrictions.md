---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_branch_restrictions"
sidebar_current: "docs-bitbucket-data-branch-restrictions"
description: |-
  Provides information about Bitbucket branch restrictions.
---

# bitbucket\_branch\_restrictions

Provides information about Bitbucket branch restrictions.

## Example Usage

```hcl
data "bitbucket_branch_restrictions" "example" {
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

* `id` - The identifier of the branch restrictions.
* `restrictions` - The restrictions. Each item contains:
    * `enabled` - The enabled.
    * `groups` - The groups.
    * `id` - The id.
    * `kind` - The kind.
    * `links` - The links.
    * `pattern` - The pattern.
    * `users` - The users.
    * `value` - The value.
