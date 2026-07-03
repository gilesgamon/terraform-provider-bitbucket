---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_issue"
sidebar_current: "docs-bitbucket-data-issue"
description: |-
  Provides information about Bitbucket issue.
---

# bitbucket\_issue

Provides information about Bitbucket issue.

## Example Usage

```hcl
data "bitbucket_issue" "example" {
  issue_id = "issue_id"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `issue_id` - (Required) The issue id.
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the issue.
* `assignee` - The assignee.
* `content` - The content.
* `created_on` - The created on.
* `kind` - The kind.
* `links` - The links.
* `priority` - The priority.
* `reporter` - The reporter.
* `state` - The state.
* `title` - The title.
* `updated_on` - The updated on.
