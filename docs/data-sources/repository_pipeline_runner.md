---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_runner"
sidebar_current: "docs-bitbucket-data-repository-pipeline-runner"
description: |-
  Provides a single self-hosted Bitbucket Pipelines runner for a repository
---

# bitbucket\_repository\_pipeline\_runner

Retrieves a single self-hosted Bitbucket Pipelines runner for a repository by UUID.

OAuth2 Scopes: `runner:read`

## Example Usage

```hcl
data "bitbucket_repository_pipeline_runner" "example" {
  workspace   = "gob"
  repo_slug   = "example"
  runner_uuid = "{12345678-90ab-cdef-1234-567890abcdef}"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.
* `runner_uuid` - (Required) The UUID of the runner surrounded by curly-braces.

## Attributes Reference

* `uuid` - The UUID identifying the runner.
* `name` - The name of the runner.
* `labels` - Labels assigned to the runner.
* `created_on` - The timestamp when the runner was created.
* `updated_on` - The timestamp when the runner was last updated.
* `state` - The runner state block containing `status`, `cordoned`, `updated_on` and a `version` map.
* `oauth_client` - The OAuth client configuration for runner authentication.
