---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_runners"
sidebar_current: "docs-bitbucket-data-repository-pipeline-runners"
description: |-
  Provides a list of self-hosted Bitbucket Pipelines runners for a repository
---

# bitbucket\_repository\_pipeline\_runners

Retrieves the list of self-hosted Bitbucket Pipelines runners configured for a repository.

OAuth2 Scopes: `runner:read`

## Example Usage

```hcl
data "bitbucket_repository_pipeline_runners" "example" {
  workspace = "gob"
  repo_slug = "example"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.

## Attributes Reference

* `id` - The identifier of the runner collection.
* `runners` - A list of runners. See the [workspace runners](workspace_pipeline_runners.md#runner) data source for the nested attributes of each runner.
