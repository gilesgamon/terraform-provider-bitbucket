---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_team_pipeline_variables"
sidebar_current: "docs-bitbucket-data-team-pipeline-variables"
description: |-
  Provides information about Bitbucket team pipeline variables.
---

# bitbucket\_team\_pipeline\_variables

Provides information about Bitbucket team pipeline variables.

## Example Usage

```hcl
data "bitbucket_team_pipeline_variables" "example" {
  username = "username"
}
```

## Argument Reference

The following arguments are supported:

* `username` - (Required) Team username

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the team pipeline variables.
* `variables` - The variables. Each item contains:
    * `created_on` - Creation timestamp
    * `key` - Variable key
    * `secured` - Whether the variable is secured
    * `updated_on` - Last update timestamp
    * `uuid` - Variable UUID
    * `value` - Variable value
