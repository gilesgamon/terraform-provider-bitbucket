## 0.2.0 (July 2026)

### 🔐 Provider configuration

* Added descriptions to all provider configuration attributes (improves the registry documentation).
* Marked `password`, `oauth_client_secret` and `oauth_token` as `Sensitive` so they are redacted in Terraform output.
* Added `make lint` and `make docs` targets.

### ⚡ Pagination

* List/collection data sources now return the **full** result set instead of only the first page (Bitbucket defaults to 10 items per page). A `Client.GetAll` helper transparently follows `next` links, and 59 collection data sources were switched to it. Endpoints that take an explicit `page` argument (e.g. code search) are intentionally left page-by-page.

### ⚡ Client

* The API client now sends a `User-Agent` header (`terraform-provider-bitbucket/<version>`), with the version injected at build time.
* Requests that are rate limited (HTTP 429) are now retried with backoff, honouring the `Retry-After` header when present.

### 🧪 CI & Testing

* Added a `Test` GitHub Actions workflow that runs `gofmt`, `go vet`, `go build` and `go test` on every pull request and on pushes to `main`.
* Added credential-free unit tests for pure helpers (query encoding, pagination link handling, ID parsers, and flatten functions).

### 📦 Project

* Rewrote the README with accurate resource/data-source counts, authentication guidance and examples.
* Added `CONTRIBUTING.md`, `ROADMAP.md`, issue templates and a pull request template.
* Added runnable configurations under `examples/`, a `terraform-registry-manifest.json`, and a `.go-version` file.
* Removed internal work-log documents from the repository root.

### ✨ New Data Sources

* `bitbucket_tags` - List all Git tags in a repository (complements the singular `bitbucket_tag`).
* `bitbucket_effective_branching_model` - The effective branching model for a repository, including project-level inheritance.
* `bitbucket_branch_merge_base` - The common ancestor commit between two revisions.

### 📖 Documentation

* Added reference documentation for every registered resource and data source, closing a gap where 88 data sources and 2 resources (including the newer pipeline runner, conflicts, user-workspace, and project deploy key endpoints) had no docs page. List/collection data sources now document their nested item attributes.
* Added a schema-driven documentation generator under `tools/docgen` to keep `docs/` in sync with the provider schema for future endpoints.
* Fixed a typo on the provider index page.

### 🧹 Cleanup

* Removed 73 unused data source implementations that were never registered with the provider. These were either duplicates of existing data sources or targeted endpoints that do not exist in the Bitbucket Cloud API (for example, a fabricated `issue-fields/*` tree, per-step pipeline sub-resources such as `steps/{uuid}/max-seconds`, and an `addons/*/webhooks/*` tree). This removes dead code and the associated linter noise.
* Fixed the endpoint URLs for the merge base (`/merge-base/{revspec}`) and effective branching model (`/effective-branching-model`) data sources so they match the real API, and corrected the merge base response to return the ancestor commit.
* Formatted all Go source with `gofmt` and fixed an error-string lint warning.

## 0.1.6 (July 2026)

### ✨ New Features

* **Project Deploy Keys** (resolves #1): Added the `bitbucket_project_deploy_key` resource and the `bitbucket_project_deploy_keys` data source, covering the previously unimplemented `/workspaces/{workspace}/projects/{project_key}/deploy-keys` endpoints. Project deploy keys are inherited by all repositories in a project.

### ✨ New Endpoints (latest Bitbucket Cloud API sync)

Synced the vendored OpenAPI reference (`reference/swagger.v3.json`) with the
latest Bitbucket Cloud API and implemented the newly added endpoints.

#### New Data Sources

* **Pipeline Runners**: `bitbucket_workspace_pipeline_runners`, `bitbucket_workspace_pipeline_runner`, `bitbucket_repository_pipeline_runners`, `bitbucket_repository_pipeline_runner`
* **Merge Conflicts**: `bitbucket_file_conflicts`, `bitbucket_pull_request_conflicts`
* **User Workspace Access**: `bitbucket_user_workspaces`, `bitbucket_user_workspace_permission`, `bitbucket_user_workspace_repository_permissions`
* **Connect Add-on**: `bitbucket_addon_client_key`
* **Workspace GPG**: `bitbucket_workspace_gpg_public_key`

#### New Resources

* **Pipeline Runners**: `bitbucket_workspace_pipeline_runner`, `bitbucket_repository_pipeline_runner` — manage self-hosted Bitbucket Pipelines runners at the workspace and repository level.

### 🗑️ Deprecated Upstream Endpoints

The following `addon/linkers` endpoints were removed from the upstream Bitbucket
Cloud API. The `bitbucket_repository_addon_linkers` and
`bitbucket_repository_addon_values` data sources are retained for backwards
compatibility but now target endpoints that Atlassian no longer documents:

* `GET /addon/linkers`
* `GET /addon/linkers/{linker_key}`
* `GET|POST|PUT|DELETE /addon/linkers/{linker_key}/values`
* `GET|DELETE /addon/linkers/{linker_key}/values/{value_id}`

### 🐛 Bug Fixes

* **HTTP client crash**: Fixed a nil-pointer dereference (segfault) in the HTTP client that occurred on transport-level errors (DNS/connection failures/timeouts). The client now returns the underlying error instead of panicking when no response is received.
* **`bitbucket_workspace_members`**: Fixed swapped `username` and `display_name` attributes.

### ⚡ Improvements

* **Pagination**: Added a `GetPaginated` client helper that follows Bitbucket `next` links, and adopted it in the `bitbucket_users`, `bitbucket_groups`, `bitbucket_projects`, and `bitbucket_workspaces` data sources so they return the full result set instead of only the first page (default page size of 10).
* **Query encoding**: Query parameters (for example the `q` filter) are now URL-encoded and emitted in a deterministic order, preventing malformed requests when values contain spaces or special characters.

## 2.0.0 (December 2024)

### 🎉 Major Release - Complete API Coverage

This release represents a complete overhaul of the Bitbucket Terraform Provider with **100% API coverage** based on the latest Bitbucket API v3 specification.

### ✨ New Features

#### **Complete API Implementation (178 endpoints)**
* **Data Sources**: 86 new data sources covering all Bitbucket API endpoints
* **Resources**: 92 resources for comprehensive Bitbucket management
* **Total Coverage**: 178/178 endpoints (100% complete)

#### **New Data Sources**
* **Snippets Management**: `bitbucket_snippet`, `bitbucket_snippets`
* **Code Search**: `bitbucket_user_search_code`, `bitbucket_workspace_search_code`, `bitbucket_team_search_code`
* **GPG Key Management**: `bitbucket_user_gpg_key`, `bitbucket_user_gpg_keys`
* **Team Management**: `bitbucket_team_pipeline_variable`, `bitbucket_team_pipeline_variables`
* **User Email Management**: `bitbucket_user_email`, `bitbucket_user_emails`
* **Repository Settings**: `bitbucket_repository_override_settings`
* **Advanced Pipeline Features**: `bitbucket_pipeline_build_number`, `bitbucket_pipeline_schedule_executions`, `bitbucket_pipeline_test_cases`, `bitbucket_pipeline_test_case_reasons`
* **Issue Management**: `bitbucket_repository_issue_import`, `bitbucket_repository_issue_export_status`
* **Advanced PR Features**: `bitbucket_pull_request_tasks`, `bitbucket_pull_request_task`, `bitbucket_pull_request_merge_task_status`

#### **New Resources**
* **Snippet Management**: `bitbucket_snippet` - Create and manage code snippets
* **GPG Key Management**: `bitbucket_user_gpg_key` - Manage user GPG keys
* **Pipeline Control**: `bitbucket_pipeline_stop` - Stop running pipelines

#### **Enhanced Authentication**
* **OAuth 2.0 Support**: Complete OAuth client credentials flow
* **Multiple Auth Methods**: Username/password, OAuth client credentials, and OAuth tokens
* **Environment Variable Support**: All authentication methods support environment variables

### 🐛 Bug Fixes

* **Critical Fix**: Resolved nil pointer dereference in `bitbucket_repository` resource that was causing provider crashes
* **Schema Validation**: Fixed `bitbucket_snippet` resource ID field type from `TypeInt` to `TypeString` to comply with Terraform requirements
* **Error Handling**: Improved error handling across all resources and data sources
* **Type Safety**: Fixed type conflicts and compilation errors in new implementations

### 📚 Documentation

* **Complete API Reference**: Updated documentation for all 178 endpoints
* **Swagger Integration**: Added reference to latest Bitbucket API v3 specification (`reference/swagger.v3.json`)
* **Usage Examples**: Comprehensive examples for all new features
* **Authentication Guide**: Detailed OAuth setup and configuration examples

### 🔧 Technical Improvements

* **Code Quality**: Applied `go fmt` formatting across entire codebase
* **Linting**: Resolved all linter warnings and errors
* **Testing**: All unit tests passing with improved coverage
* **Build System**: Streamlined build process and dependency management

### 📋 Migration Notes

* **Breaking Changes**: This is a major version update with significant new functionality
* **Backward Compatibility**: Existing resources remain compatible
* **New Dependencies**: Updated to latest Terraform Plugin SDK v2
* **Authentication**: New OAuth options available but existing auth methods still supported

---

## 1.3.0 (March 15, 2021)

This release contains the changes to upstream repo that were never released.

### Features

* add `bitbucket_deployment` and `bitbucket_deployment_variable` resources [#60](https://github.com/hashicorp/terraform-provider-bitbucket/pull/60)
* add `require_default_reviewer_approvals_to_merge` branch restriction value [#52](https://github.com/hashicorp/terraform-provider-bitbucket/pull/52)

### Bug fixes

* fix issue with omitempty [#49](https://github.com/hashicorp/terraform-provider-bitbucket/pull/49)

### Documentation

* fix ducmentation typo [#54](https://github.com/hashicorp/terraform-provider-bitbucket/pull/54), [#61](https://github.com/hashicorp/terraform-provider-bitbucket/pull/61) and [#65](https://github.com/hashicorp/terraform-provider-bitbucket/pull/65)

## 1.2.0 (January 23, 2020)
* add `bitbucket_project` to create a new project via the API
* add `bitbucket_repository` turn on/off pipelines
* add `bitbucket_repository_variable` to add variables via terraform to your pipelines builds
* add `bitbucket_user` to find a user and use for default reviewers.

## 1.1.0 (June 19, 2019)

### Features

* add `skip_cert_verification` for hooks [#19](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/19)

### Bug fixes

* handle missing hooks [#24](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/24)
* fix default reviewer pagination bug [#28](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/28)

### Dev updates

* add `website` and `website-test` targets [#16](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/16)
* add `website-test` target to Travis [#17](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/17)
* upgrade to go 1.11 [#25](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/25)
* switch to go modules [#27](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/27)
* upgrade to `hashicorp/terraform` v0.12.2 [#34](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/34)

### Documentation

* add note about v1 APIs [#21](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/21)

## 1.0.0 (December 08, 2017)

* resource/bitbucket_repository: Add the ability to define a seperate slug for a repository ([#5](https://github.com/terraform-providers/terraform-provider-bitbucket/issues/5))

## 0.1.0 (June 20, 2017)

NOTES:

* Same functionality as that of Terraform 0.9.8. Repacked as part of [Provider Splitout](https://www.hashicorp.com/blog/upcoming-provider-changes-in-terraform-0-10/)
