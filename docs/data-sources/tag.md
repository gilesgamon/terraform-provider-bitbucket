---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_tag"
sidebar_current: "docs-bitbucket-data-tag"
description: |-
  Provides information about Bitbucket tag.
---

# bitbucket\_tag

Provides information about Bitbucket tag.

## Example Usage

```hcl
data "bitbucket_tag" "example" {
  repo_slug = "example-repo"
  tag_name = "tag_name"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) The repo slug.
* `tag_name` - (Required) The tag name.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the tag.
* `author` - The author. Each item contains:
    * `display_name` - The display name.
    * `username` - The username.
    * `uuid` - The uuid.
* `message` - The message.
* `name` - The name.
* `target_date` - The target date.
* `target_hash` - The target hash.
* `uuid` - The uuid.
