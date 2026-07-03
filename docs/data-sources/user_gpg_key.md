---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_gpg_key"
sidebar_current: "docs-bitbucket-data-user-gpg-key"
description: |-
  Provides information about Bitbucket user gpg key.
---

# bitbucket\_user\_gpg\_key

Provides information about Bitbucket user gpg key.

## Example Usage

```hcl
data "bitbucket_user_gpg_key" "example" {
  fingerprint = "fingerprint"
  selected_user = "selected_user"
}
```

## Argument Reference

The following arguments are supported:

* `fingerprint` - (Required) GPG key fingerprint
* `selected_user` - (Required) User UUID or username

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the user gpg key.
* `created_on` - Creation timestamp
* `key` - GPG public key content
* `links` - GPG key links Each item contains:
    * `self` - The self.
* `owner` - Key owner Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `type` - GPG key type
