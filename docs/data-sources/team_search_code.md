---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_team_search_code"
sidebar_current: "docs-bitbucket-data-team-search-code"
description: |-
  Provides information about Bitbucket team search code.
---

# bitbucket\_team\_search\_code

Provides information about Bitbucket team search code.

## Example Usage

```hcl
data "bitbucket_team_search_code" "example" {
  search_query = "search_query"
  username = "username"
}
```

## Argument Reference

The following arguments are supported:

* `search_query` - (Required) The search query string
* `username` - (Required) Team username
* `page` - (Optional) Page number for pagination
* `pagelen` - (Optional) Number of results per page

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `id` - The identifier of the team search code.
* `query_substituted` - The actual query that was executed
* `results` - The results. Each item contains:
    * `content_match_count` - Number of content matches
    * `content_matches` - Content matches
    * `file` - File information
    * `path_matches` - Path matches
    * `type` - Result type
