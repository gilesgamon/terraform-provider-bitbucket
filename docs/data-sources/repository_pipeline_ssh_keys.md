---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_ssh_keys"
sidebar_current: "docs-bitbucket-data-repository-pipeline-ssh-keys"
description: |-
  Provides information about Bitbucket repository pipeline ssh keys.
---

# bitbucket\_repository\_pipeline\_ssh\_keys

Provides information about Bitbucket repository pipeline ssh keys.

## Example Usage

```hcl
data "bitbucket_repository_pipeline_ssh_keys" "example" {
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

* `id` - The identifier of the repository pipeline ssh keys.
* `ssh_keys` - The ssh keys. Each item contains:
    * `comment` - The comment.
    * `created_on` - The created on.
    * `label` - The label.
    * `links` - The links.
    * `public_key` - The public key.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
