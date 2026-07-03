---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_pull_request_conflicts"
sidebar_current: "docs-bitbucket-data-pull-request-conflicts"
description: |-
  Provides the file conflicts for a pull request
---

# bitbucket\_pull\_request\_conflicts

Retrieves the list of file conflicts for a pull request.

OAuth2 Scopes: `pullrequest`

## Example Usage

```hcl
data "bitbucket_pull_request_conflicts" "example" {
  workspace       = "gob"
  repo_slug       = "example"
  pull_request_id = "42"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) This can either be the workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.
* `pull_request_id` - (Required) The ID of the pull request.

## Attributes Reference

* `id` - The identifier of the conflict collection.
* `conflicts` - A list of file conflicts. See the [file conflicts](file_conflicts.md#file-conflict) data source for the nested attributes.
