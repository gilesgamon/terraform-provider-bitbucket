---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_project_deploy_key"
sidebar_current: "docs-bitbucket-resource-project-deploy-key"
description: |-
  Provides a Bitbucket Project Deploy Key
---

# bitbucket\_project\_deploy\_key

Provides a Bitbucket Project Deploy Key resource.

This allows you to manage deploy (access) keys at the project level. Project
deploy keys are inherited by all repositories in the project, including
repositories created outside of Terraform, which is useful for granting a shared
read-only key to every repository under a project.

OAuth2 Scopes: `project` and `project:admin`

## Example Usage

```hcl
resource "bitbucket_project_deploy_key" "example" {
  workspace   = "example"
  project_key = "PROJ"
  key         = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIKqP3Cr632C2dNhhgKVcon4ldUSAeKiku2yP9O9/bDtY"
  label       = "shared-read-key"
}
```

## Argument Reference

The following arguments are supported:

* `workspace` - (Required) The workspace ID (slug) or the workspace UUID surrounded by curly-braces. Changing this forces a new resource.
* `project_key` - (Required) The project key (for example `PROJ`). Changing this forces a new resource.
* `key` - (Required) The SSH public key value in OpenSSH format. Changing this forces a new resource.
* `label` - (Optional) The user-defined label for the deploy key. Changing this forces a new resource.

## Attributes Reference

In addition to the arguments above, the following attributes are exported:

* `key_id` - The deploy key's ID.
* `comment` - The comment parsed from the deploy key (if present).
* `added_on` - The timestamp when the deploy key was added.
* `last_used` - The timestamp when the deploy key was last used.

## Import

Project deploy keys can be imported using their `workspace/project-key/key-id` ID, e.g.

```sh
terraform import bitbucket_project_deploy_key.example workspace/PROJ/1234
```
