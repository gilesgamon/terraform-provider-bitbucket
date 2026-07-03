---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_team_pipeline_variable"
sidebar_current: "docs-bitbucket-data-team-pipeline-variable"
description: |-
  Provides information about Bitbucket team pipeline variable.
---

# bitbucket\_team\_pipeline\_variable

Provides information about Bitbucket team pipeline variable.

## Example Usage

```hcl
data "bitbucket_team_pipeline_variable" "example" {
  username = "username"
  variable_uuid = "variable_uuid"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) Team username
* `variable_uuid` - (Required) Variable UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the team pipeline variable.
* `created_on` - Creation timestamp
* `key` - Variable key
* `secured` - Whether the variable is secured
* `updated_on` - Last update timestamp
* `uuid` - Variable UUID
* `value` - Variable value
