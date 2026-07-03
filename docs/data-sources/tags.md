---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_tags"
sidebar_current: "docs-bitbucket-data-tags"
description: |-
  Provides the list of tags for a Bitbucket repository
---

# bitbucket\_tags

Retrieves the list of Git tags for a repository. Use the singular
`bitbucket_tag` data source to look up a single tag by name.

OAuth2 Scopes: `repository`

## Example Usage

```hcl
data "bitbucket_tags" "example" {
  workspace = "example-workspace"
  repo_slug = "example-repo"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `repo_slug` - (Required) The repository slug.

## Attributes Reference

* `id` - The identifier of the tags collection.
* `tags` - A list of tags. See [Tag](#tag) below.

### Tag

* `name` - The name of the tag.
* `target` - The commit the tag points to.
* `message` - The tag message (for annotated tags).
* `tagger` - The tagger of the tag (for annotated tags).
* `date` - The date the tag was created.
* `links` - Links related to the tag.
