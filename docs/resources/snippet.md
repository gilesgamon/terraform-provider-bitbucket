# bitbucket_snippet

Provides a Bitbucket snippet resource.

This allows snippets to be created, read, updated, and deleted.

## Example Usage

```hcl
resource "bitbucket_snippet" "example" {
  workspace = "my-workspace"
  title     = "Example Snippet"
  is_private = true
  
  files = {
    "main.py" = <<-EOT
      def hello_world():
          print("Hello, World!")
      
      if __name__ == "__main__":
          hello_world()
    EOT
    
    "README.md" = <<-EOT
      # Example Snippet
      
      This is an example Python snippet.
    EOT
  }
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace slug or UUID.
* `title` - (Required) The snippet title.
* `files` - (Required) A map of filename to file content.
* `is_private` - (Optional) Whether the snippet is private. Defaults to `true`.
* `scm` - (Optional) The DVCS used to store the snippet. Defaults to `git`.

## Attributes Reference

In addition to all arguments above, the following attributes are exported:

* `id` - The snippet ID in the format `workspace/encoded_id`.
* `snippet_id` - The numeric snippet ID.
* `encoded_id` - The snippet encoded ID.
* `created_on` - The creation timestamp.
* `updated_on` - The last update timestamp.
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

## Import

Snippets can be imported using the `workspace/encoded_id` format:

```
$ terraform import bitbucket_snippet.example my-workspace/abc123def456
```

