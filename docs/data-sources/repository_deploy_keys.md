---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_deploy_keys"
sidebar_current: "docs-bitbucket-data-repository-deploy-keys"
description: |-
  Provides information about Bitbucket repository deploy keys.
---

# bitbucket\_repository\_deploy\_keys

Provides information about Bitbucket repository deploy keys.

## Example Usage

```hcl
data "bitbucket_repository_deploy_keys" "example" {
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

* `id` - The identifier of the repository deploy keys.
* `deploy_keys` - The deploy keys. Each item contains:
    * `comment` - The comment.
    * `created_on` - The created on.
    * `id` - The id.
    * `key` - The key.
    * `label` - The label.
    * `last_used` - The last used.
    * `links` - The links.
