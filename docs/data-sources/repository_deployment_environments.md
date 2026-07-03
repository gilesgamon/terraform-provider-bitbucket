---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_deployment_environments"
sidebar_current: "docs-bitbucket-data-repository-deployment-environments"
description: |-
  Provides information about Bitbucket repository deployment environments.
---

# bitbucket\_repository\_deployment\_environments

Provides information about Bitbucket repository deployment environments.

## Example Usage

```hcl
data "bitbucket_repository_deployment_environments" "example" {
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

* `id` - The identifier of the repository deployment environments.
* `environments` - The environments. Each item contains:
    * `created_on` - The created on.
    * `deployment_gate_check` - The deployment gate check.
    * `deployment_gate_enabled` - The deployment gate enabled.
    * `environment_type` - The environment type.
    * `links` - The links.
    * `name` - The name.
    * `rank` - The rank.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
