---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspace_variables"
sidebar_current: "docs-bitbucket-data-workspace-variables"
description: |-
  Provides information about Bitbucket workspace variables.
---

# bitbucket\_workspace\_variables

Provides information about Bitbucket workspace variables.

## Example Usage

```hcl
data "bitbucket_workspace_variables" "example" {
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the workspace variables.
* `variables` - The variables. Each item contains:
    * `created_on` - The created on.
    * `key` - The key.
    * `secured` - The secured.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
    * `value` - The value.
