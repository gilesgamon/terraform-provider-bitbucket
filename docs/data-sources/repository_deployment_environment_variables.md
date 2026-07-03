---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_deployment_environment_variables"
sidebar_current: "docs-bitbucket-data-repository-deployment-environment-variables"
description: |-
  Provides information about Bitbucket repository deployment environment variables.
---

# bitbucket\_repository\_deployment\_environment\_variables

Provides information about Bitbucket repository deployment environment variables.

## Example Usage

```hcl
data "bitbucket_repository_deployment_environment_variables" "example" {
  environment_uuid = "environment_uuid"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `environment_uuid` - (Required) The environment uuid.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository deployment environment variables.
* `variables` - The variables. Each item contains:
    * `created_on` - The created on.
    * `key` - The key.
    * `links` - The links.
    * `secured` - The secured.
    * `type` - The type.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
    * `value` - The value.
