---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_commit_properties"
sidebar_current: "docs-bitbucket-data-commit-properties"
description: |-
  Provides information about Bitbucket commit properties.
---

# bitbucket\_commit\_properties

Provides information about Bitbucket commit properties.

## Example Usage

```hcl
data "bitbucket_commit_properties" "example" {
  commit = "a1b2c3d4"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `commit` - (Required) The commit.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the commit properties.
* `properties` - The properties. Each item contains:
    * `key` - The key.
    * `value` - The value.
