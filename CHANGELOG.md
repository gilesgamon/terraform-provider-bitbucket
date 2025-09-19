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
