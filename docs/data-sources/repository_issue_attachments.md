---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_repository_issue_attachments"
sidebar_current: "docs-bitbucket-data-repository-issue-attachments"
description: |-
  Provides information about Bitbucket repository issue attachments.
---

# bitbucket\_repository\_issue\_attachments

Provides information about Bitbucket repository issue attachments.

## Example Usage

```hcl
data "bitbucket_repository_issue_attachments" "example" {
  issue_id = "issue_id"
  repo_slug = "example-repo"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `issue_id` - (Required) Issue ID
* `repo_slug` - (Required) The repo slug.
* `workspace` - (Required) The workspace.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the repository issue attachments.
* `attachments` - The attachments. Each item contains:
    * `created_on` - The created on.
    * `links` - The links.
    * `name` - The name.
    * `size` - The size.
