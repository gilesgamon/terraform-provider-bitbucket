# Roadmap

Planned and in-progress improvements. This is a living document — items may be
reprioritised.

## In progress

- **Full pagination coverage.** A reusable `Client.GetPaginated` helper exists
  and is used by the core collection data sources. It is being rolled out to the
  remaining list data sources so they return the complete result set rather than
  only the first page.
- **Correct nested-object attributes.** Some data sources model nested API
  objects (`links`, `target`, `owner`, ...) as `TypeMap` of strings, which
  silently drops nested values. These are being migrated to typed nested blocks
  or JSON string attributes.

## Planned

- **Complete `tfplugindocs` migration.** Documentation is complete and generated
  via `tools/docgen`, and `examples/` + a registry manifest are in place. The
  remaining work is to add `Description`s to every schema attribute and switch
  documentation generation to
  [`tfplugindocs`](https://github.com/hashicorp/terraform-plugin-docs), with a CI
  check that fails when generated docs are stale. Field descriptions are being
  filled in incrementally.
- **Shared read helpers.** Extract the repeated get→read→unmarshal→not-found
  boilerplate in data sources into shared helpers to reduce duplication and
  centralise error handling.
- **Client resilience.** Automatic retry with backoff on HTTP 429 (Bitbucket
  rate limiting).

## Long term

- **Migrate to `terraform-plugin-framework`.** The provider currently uses
  [`terraform-plugin-sdk/v2`](https://github.com/hashicorp/terraform-plugin-sdk).
  Migrating to the newer
  [`terraform-plugin-framework`](https://github.com/hashicorp/terraform-plugin-framework)
  would bring first-class support for nested attributes, plan modifiers, and
  richer validation. This is a large, incremental effort (it can run alongside
  SDKv2 via `terraform-plugin-mux`) and is not yet scheduled.
