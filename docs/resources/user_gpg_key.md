---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_user_gpg_key"
sidebar_current: "docs-bitbucket-resource-user-gpg-key"
description: |-
  Provides a Bitbucket user gpg key resource.
---

# bitbucket\_user\_gpg\_key

Provides a Bitbucket user gpg key resource.

## Example Usage

```hcl
resource "bitbucket_user_gpg_key" "example" {
  key = "key"
  selected_user = "selected_user"
}
```

## Argument Reference

The following arguments are supported:

* `key` - (Required) GPG public key content
* `selected_user` - (Required) User UUID or username

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the user gpg key.
* `created_on` - Creation timestamp
* `fingerprint` - GPG key fingerprint
* `links` - GPG key links Each item contains:
    * `self` - The self.
* `owner` - Key owner Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `type` - GPG key type
