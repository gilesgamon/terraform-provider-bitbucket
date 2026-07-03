---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_effective_branching_model"
sidebar_current: "docs-bitbucket-data-effective-branching-model"
description: |-
  Provides the effective branching model for a Bitbucket repository
---

# bitbucket\_effective\_branching\_model

Retrieves the *effective* branching model for a repository — the branching
model that is actually applied, taking any project-level inheritance into
account. Use the `bitbucket_branching_model` data source for the repository's
own configured model.

OAuth2 Scopes: `repository`

## Example Usage

```hcl
data "bitbucket_effective_branching_model" "example" {
  workspace = "example-workspace"
  repo_slug = "example-repo"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.

## Attributes Reference

* `id` - The identifier of the effective branching model.
* `development` - The development branch configuration.
* `production` - The production branch configuration.
* `branch_types` - The configured branch types. Each item contains `kind`, `prefix` and `enabled`.
* `links` - Links related to the branching model.
