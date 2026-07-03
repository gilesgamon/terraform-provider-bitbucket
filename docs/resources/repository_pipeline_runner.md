---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_runner"
sidebar_current: "docs-bitbucket-resource-repository-pipeline-runner"
description: |-
  Manages a self-hosted Bitbucket Pipelines runner for a repository
---

# bitbucket\_repository\_pipeline\_runner

Manages a self-hosted Bitbucket Pipelines runner at the repository level.

OAuth2 Scopes: `runner:write`

## Example Usage

```hcl
resource "bitbucket_repository_pipeline_runner" "example" {
  workspace = "gob"
  repo_slug = "example"
  name      = "linux-runner"
  labels    = ["self.hosted", "linux", "shell"]
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces. Changing this forces a new resource.
* `repo_slug` - (Required) The repository slug. Changing this forces a new resource.
* `name` - (Required) The name of the runner.
* `labels` - (Required) The set of labels assigned to the runner for identification and routing.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the runner in the form `WORKSPACE/REPO-SLUG/RUNNER-UUID`.
* `uuid` - The UUID identifying the runner.
* `state` - A map describing the runner state (`status`, `cordoned`, `updated_on`).
* `oauth_client` - The OAuth client configuration for runner authentication. Marked sensitive; the `secret` value is only returned once when the runner is created.
* `created_on` - The timestamp when the runner was created.
* `updated_on` - The timestamp when the runner was last updated.

## Import

Repository pipeline runners can be imported using the workspace, repository slug and runner UUID:

```sh
terraform import bitbucket_repository_pipeline_runner.example gob/example/{12345678-90ab-cdef-1234-567890abcdef}
```
