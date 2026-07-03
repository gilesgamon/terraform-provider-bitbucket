---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_issues"
sidebar_current: "docs-bitbucket-data-issues"
description: |-
  Provides information about Bitbucket issues.
---

# bitbucket\_issues

Provides information about Bitbucket issues.

## Example Usage

```hcl
data "bitbucket_issues" "example" {
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

* `id` - The identifier of the issues.
* `issues` - The issues. Each item contains:
    * `assignee` - The assignee.
    * `content` - The content.
    * `created_on` - The created on.
    * `id` - The id.
    * `kind` - The kind.
    * `links` - The links.
    * `priority` - The priority.
    * `reporter` - The reporter.
    * `state` - The state.
    * `title` - The title.
    * `updated_on` - The updated on.
