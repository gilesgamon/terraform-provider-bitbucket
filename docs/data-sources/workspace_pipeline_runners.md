---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspace_pipeline_runners"
sidebar_current: "docs-bitbucket-data-workspace-pipeline-runners"
description: |-
  Provides a list of self-hosted Bitbucket Pipelines runners for a workspace
---

# bitbucket\_workspace\_pipeline\_runners

Retrieves the list of self-hosted Bitbucket Pipelines runners configured for a workspace.

OAuth2 Scopes: `runner:read`

## Example Usage

```hcl
data "bitbucket_workspace_pipeline_runners" "example" {
  workspace = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.

## Attributes Reference

* `id` - The identifier of the runner collection.
* `runners` - A list of runners. See [Runner](#runner) below.

### Runner

* `uuid` - The UUID identifying the runner.
* `name` - The name of the runner.
* `labels` - Labels assigned to the runner for identification and routing.
* `created_on` - The timestamp when the runner was created.
* `updated_on` - The timestamp when the runner was last updated.
* `state` - The runner state block containing `status`, `cordoned`, `updated_on` and a `version` map.
* `oauth_client` - The OAuth client configuration for runner authentication (`id`, `token_endpoint`, `audience`). The `secret` is only returned when the runner is created.
