# bitbucket_snippets

Use this data source to access information about Bitbucket snippets.

## Example Usage

```hcl
# Get all snippets for a workspace
data "bitbucket_snippets" "workspace_snippets" {
  workspace = "my-workspace"
}

# Get all snippets (global)
data "bitbucket_snippets" "all_snippets" {
}

output "snippet_count" {
  value = length(data.bitbucket_snippets.workspace_snippets.snippets)
}

output "first_snippet_title" {
  value = data.bitbucket_snippets.workspace_snippets.snippets[0].title
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Optional) The workspace slug or UUID to filter snippets. If not provided, returns all snippets.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The data source ID.
* `snippets` - List of snippets.
  * `id` - The numeric snippet ID.
  * `title` - The snippet title.
  * `scm` - The DVCS used to store the snippet.
  * `created_on` - The creation timestamp.
  * `updated_on` - The last update timestamp.
  * `is_private` - Whether the snippet is private.
  * `owner` - The snippet owner information.
    * `username` - Owner username.
    * `display_name` - Owner display name.
    * `uuid` - Owner UUID.
  * `creator` - The snippet creator information.
    * `username` - Creator username.
    * `display_name` - Creator display name.
    * `uuid` - Creator UUID.
  * `links` - Snippet links.
    * `self` - Self link.
      * `href` - The URL.
    * `html` - HTML link.
      * `href` - The URL.

