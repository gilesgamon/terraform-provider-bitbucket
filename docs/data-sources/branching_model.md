---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_branching_model"
sidebar_current: "docs-bitbucket-data-branching-model"
description: |-
  Provides information about Bitbucket branching model.
---

# bitbucket\_branching\_model

Provides information about Bitbucket branching model.

## Example Usage

```hcl
data "bitbucket_branching_model" "example" {
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

* `id` - The identifier of the branching model.
* `branch_types` - The branch types. Each item contains:
    * `enabled` - The enabled.
    * `kind` - The kind.
    * `prefix` - The prefix.
* `development` - The development.
* `links` - The links.
* `production` - The production.
