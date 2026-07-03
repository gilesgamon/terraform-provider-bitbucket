---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_ssh_keys"
sidebar_current: "docs-bitbucket-data-ssh-keys"
description: |-
  Provides information about Bitbucket ssh keys.
---

# bitbucket\_ssh\_keys

Provides information about Bitbucket ssh keys.

## Example Usage

```hcl
data "bitbucket_ssh_keys" "example" {
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the ssh keys.
* `ssh_keys` - The ssh keys. Each item contains:
    * `comment` - The comment.
    * `created_on` - The created on.
    * `key` - The key.
    * `label` - The label.
    * `last_used` - The last used.
    * `links` - The links.
    * `owner` - The owner.
    * `uuid` - The uuid.
