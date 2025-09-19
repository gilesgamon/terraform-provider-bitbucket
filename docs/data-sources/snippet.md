# bitbucket_snippet

Use this data source to access information about a specific Bitbucket snippet.

## Example Usage

```hcl
data "bitbucket_snippet" "example" {
  workspace   = "my-workspace"
  encoded_id  = "abc123def456"
}

output "snippet_title" {
  value = data.bitbucket_snippet.example.title
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace slug or UUID.
* `encoded_id` - (Required) The snippet encoded ID.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The snippet ID in the format `workspace/encoded_id`.
* `snippet_id` - The numeric snippet ID.
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

