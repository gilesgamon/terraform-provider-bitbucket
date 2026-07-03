---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_ssh_known_hosts"
sidebar_current: "docs-bitbucket-data-repository-pipeline-ssh-known-hosts"
description: |-
  Provides information about Bitbucket repository pipeline ssh known hosts.
---

# bitbucket\_repository\_pipeline\_ssh\_known\_hosts

Provides information about Bitbucket repository pipeline ssh known hosts.

## Example Usage

```hcl
data "bitbucket_repository_pipeline_ssh_known_hosts" "example" {
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

* `id` - The identifier of the repository pipeline ssh known hosts.
* `known_hosts` - The known hosts. Each item contains:
    * `comment` - The comment.
    * `created_on` - The created on.
    * `hostname` - The hostname.
    * `links` - The links.
    * `public_key` - The public key.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
