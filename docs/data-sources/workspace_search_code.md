---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_workspace_search_code"
sidebar_current: "docs-bitbucket-data-workspace-search-code"
description: |-
  Provides information about Bitbucket workspace search code.
---

# bitbucket\_workspace\_search\_code

Provides information about Bitbucket workspace search code.

## Example Usage

```hcl
data "bitbucket_workspace_search_code" "example" {
  search_query = "search_query"
  workspace = "example-workspace"
}
```

## Argument Reference

The following arguments are supported:

* `search_query` - (Required) The search query string
* `workspace` - (Required) Workspace slug or UUID
* `page` - (Optional) Page number for pagination
* `pagelen` - (Optional) Number of results per page

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the workspace search code.
* `query_substituted` - The actual query that was executed
* `results` - The results. Each item contains:
    * `content_match_count` - Number of content matches
    * `content_matches` - Content matches
    * `file` - File information
    * `path_matches` - Path matches
    * `type` - Result type
