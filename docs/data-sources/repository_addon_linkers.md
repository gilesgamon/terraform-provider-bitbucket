---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_addon_linkers"
sidebar_current: "docs-bitbucket-data-repository-addon-linkers"
description: |-
  Provides information about Bitbucket repository addon linkers.
---

# bitbucket\_repository\_addon\_linkers

Provides information about Bitbucket repository addon linkers.

## Example Usage

```hcl
data "bitbucket_repository_addon_linkers" "example" {
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

* `id` - The identifier of the repository addon linkers.
* `linkers` - The linkers. Each item contains:
    * `application` - The application.
    * `description` - The description.
    * `id` - The id.
    * `key` - The key.
    * `links` - The links.
    * `name` - The name.
    * `uuid` - The uuid.
    * `vendor` - The vendor.
