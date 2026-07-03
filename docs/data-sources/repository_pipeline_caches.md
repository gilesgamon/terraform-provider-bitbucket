---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_pipeline_caches"
sidebar_current: "docs-bitbucket-data-repository-pipeline-caches"
description: |-
  Provides information about Bitbucket repository pipeline caches.
---

# bitbucket\_repository\_pipeline\_caches

Provides information about Bitbucket repository pipeline caches.

## Example Usage

```hcl
data "bitbucket_repository_pipeline_caches" "example" {
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

* `id` - The identifier of the repository pipeline caches.
* `caches` - The caches. Each item contains:
    * `created_on` - The created on.
    * `last_accessed` - The last accessed.
    * `links` - The links.
    * `name` - The name.
    * `path` - The path.
    * `size` - The size.
    * `updated_on` - The updated on.
    * `uuid` - The uuid.
