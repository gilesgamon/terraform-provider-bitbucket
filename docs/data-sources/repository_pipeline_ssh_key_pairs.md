---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_ssh_key_pairs"
sidebar_current: "docs-bitbucket-data-repository-pipeline-ssh-key-pairs"
description: |-
  Provides information about Bitbucket repository pipeline ssh key pairs.
---

# bitbucket\_repository\_pipeline\_ssh\_key\_pairs

Provides information about Bitbucket repository pipeline ssh key pairs.

## Example Usage

```hcl
data "bitbucket_repository_pipeline_ssh_key_pairs" "example" {
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

* `id` - The identifier of the repository pipeline ssh key pairs.
* `key_pairs` - The key pairs. Each item contains:
    * `comment` - The comment.
    * `created_on` - The created on.
    * `label` - The label.
    * `links` - The links.
    * `private_key` - The private key.
    * `public_key` - The public key.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
