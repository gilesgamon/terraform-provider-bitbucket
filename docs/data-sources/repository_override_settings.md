---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_override_settings"
sidebar_current: "docs-bitbucket-data-repository-override-settings"
description: |-
  Provides information about Bitbucket repository override settings.
---

# bitbucket\_repository\_override\_settings

Provides information about Bitbucket repository override settings.

## Example Usage

```hcl
data "bitbucket_repository_override_settings" "example" {
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `repo_slug` - (Required) Repository slug or UUID
* `workspace` - (Required) Workspace slug or UUID

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository override settings.
* `settings` - The settings. Each item contains:
    * `links` - Settings links
    * `name` - Settings name
    * `type` - Settings type
    * `value` - Settings value
