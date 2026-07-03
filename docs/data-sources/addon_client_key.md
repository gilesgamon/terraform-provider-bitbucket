---
layout: "bitbucket"
page_title: "Bitbucket: bitbucket_addon_client_key"
sidebar_current: "docs-bitbucket-data-addon-client-key"
description: |-
  Provides the client key of a Connect add-on
---

# bitbucket\_addon\_client\_key

Retrieves the client key of the Connect add-on linked to the Forge app
installation where the request was made.

## Example Usage

```hcl
data "bitbucket_addon_client_key" "example" {
  addon_key = "my-connect-addon"
}
```

## Argument Reference

The following arguments are supported:

* `addon_key` - (Required) The key of the Connect add-on.

## Attributes Reference

* `id` - The identifier of the resource.
* `client_key` - The client key of the Connect add-on.
* `content` - The raw response body returned by the API.
