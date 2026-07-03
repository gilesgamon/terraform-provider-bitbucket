---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_project_deploy_keys"
sidebar_current: "docs-bitbucket-data-project-deploy-keys"
description: |-
  Provides the list of deploy keys configured for a Bitbucket project
---

# bitbucket\_project\_deploy\_keys

Retrieves the list of deploy (access) keys configured at the project level.

OAuth2 Scopes: `project`

## Example Usage

```hcl
data "bitbucket_project_deploy_keys" "example" {
  workspace   = "example"
  project_key = "PROJ"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace ID (slug) or the workspace UUID surrounded by curly-braces.
* `project_key` - (Required) The project key (for example `PROJ`).

## Attributes Reference

* `id` - The identifier of the deploy keys collection.
* `deploy_keys` - A list of project deploy keys. See [Deploy Key](#deploy-key) below.

### Deploy Key

* `key_id` - The deploy key's ID.
* `key` - The public SSH key value.
* `label` - The user-defined label for the deploy key.
* `comment` - The comment parsed from the deploy key (if present).
* `added_on` - The timestamp when the deploy key was added.
* `last_used` - The timestamp when the deploy key was last used.
