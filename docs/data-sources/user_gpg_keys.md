---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_gpg_keys"
sidebar_current: "docs-bitbucket-data-user-gpg-keys"
description: |-
  Provides information about Bitbucket user gpg keys.
---

# bitbucket\_user\_gpg\_keys

Provides information about Bitbucket user gpg keys.

## Example Usage

```hcl
data "bitbucket_user_gpg_keys" "example" {
  selected_user = "selected_user"
}
```

## Argument Reference

The following arguments are supported:

* `selected_user` - (Required) User UUID or username

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the user gpg keys.
* `gpg_keys` - The gpg keys. Each item contains:
    * `created_on` - Creation timestamp
    * `fingerprint` - GPG key fingerprint
    * `key` - GPG public key content
    * `links` - GPG key links
    * `owner` - Key owner
    * `type` - GPG key type
