---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspace_gpg_public_key"
sidebar_current: "docs-bitbucket-data-workspace-gpg-public-key"
description: |-
  Provides the workspace system GPG public key(s)
---

# bitbucket\_workspace\_gpg\_public\_key

Retrieves the system GPG public key(s) used by a workspace, for example to
verify commits signed on behalf of the workspace.

## Example Usage

```hcl
data "bitbucket_workspace_gpg_public_key" "example" {
  workspace = "gob"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.

## Attributes Reference

* `id` - The identifier of the resource.
* `public_key` - The workspace system GPG public key(s).
