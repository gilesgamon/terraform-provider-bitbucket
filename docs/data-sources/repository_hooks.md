---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_hooks"
sidebar_current: "docs-bitbucket-data-repository-hooks"
description: |-
  Provides information about Bitbucket repository hooks.
---

# bitbucket\_repository\_hooks

Provides information about Bitbucket repository hooks.

## Example Usage

```hcl
data "bitbucket_repository_hooks" "example" {
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

* `id` - The identifier of the repository hooks.
* `hooks` - The hooks. Each item contains:
    * `active` - The active.
    * `created_on` - The created on.
    * `description` - The description.
    * `events` - The events.
    * `skip_cert_verification` - The skip cert verification.
    * `updated_on` - The updated on.
    * `url` - The url.
    * `uuid` - The uuid.
