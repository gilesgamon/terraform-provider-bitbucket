---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_settings"
sidebar_current: "docs-bitbucket-data-repository-settings"
description: |-
  Provides information about Bitbucket repository settings.
---

# bitbucket\_repository\_settings

Provides information about Bitbucket repository settings.

## Example Usage

```hcl
data "bitbucket_repository_settings" "example" {
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

* `id` - The identifier of the repository settings.
* `created_on` - The created on.
* `description` - The description.
* `fork_policy` - The fork policy.
* `has_issues` - The has issues.
* `has_wiki` - The has wiki.
* `is_private` - The is private.
* `language` - The language.
* `links` - The links.
* `mainbranch` - The mainbranch.
* `name` - The name.
* `project` - The project.
* `scm` - The scm.
* `size` - The size.
* `updated_on` - The updated on.
* `website` - The website.
