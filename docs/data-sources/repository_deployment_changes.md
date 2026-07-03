---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_deployment_changes"
sidebar_current: "docs-bitbucket-data-repository-deployment-changes"
description: |-
  Provides information about Bitbucket repository deployment changes.
---

# bitbucket\_repository\_deployment\_changes

Provides information about Bitbucket repository deployment changes.

## Example Usage

```hcl
data "bitbucket_repository_deployment_changes" "example" {
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

* `id` - The identifier of the repository deployment changes.
* `changes` - The changes. Each item contains:
    * `created_on` - The created on.
    * `deployment` - The deployment.
    * `links` - The links.
    * `state` - The state.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
    * `version` - The version.
