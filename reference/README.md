# 📚 Bitbucket API Reference

## 🎯 Overview

This directory contains the complete Bitbucket API v3 specification used as the reference for implementing the Terraform Provider for Bitbucket.

## 📄 Files

### `swagger.v3.json`
- **Description**: Complete OpenAPI 3.0 specification for Bitbucket Cloud API v3
- **Size**: ~357KB (comprehensive API coverage)
- **Endpoints**: 178 total endpoints
- **Last Updated**: December 2024

## 🔍 API Coverage

The Terraform Provider for Bitbucket implements **100% of the endpoints** defined in this swagger specification:

### ✅ **Data Sources (86 endpoints)**
- Repository & Git Operations
- Pipeline & CI/CD Management
- Issue Tracking & Management
- User & Team Management
- Code Search & Snippets
- Advanced Features

### ✅ **Resources (92 endpoints)**
- Repository Management
- Project & Workspace Management
- Pipeline Configuration
- Permission & Security Management
- Advanced Resource Management

## 🛠️ Usage

### For Developers
This swagger file serves as the authoritative source for:
- API endpoint definitions
- Request/response schemas
- Authentication requirements
- Parameter specifications
- Error handling patterns

### For Users
The swagger specification helps understand:
- Available API capabilities
- Data structures and types
- Authentication methods
- Rate limiting and constraints

## 🔗 Related Documentation

- **Provider Documentation**: [README.md](../README.md)
- **Development Guide**: [DEVELOPMENT_README.md](../DEVELOPMENT_README.md)
- **Changelog**: [CHANGELOG.md](../CHANGELOG.md)
- **Terraform Registry**: [Provider Documentation](https://registry.terraform.io/providers/gilesgamon/terraform-provider-bitbucket/latest/docs)

## 📊 Implementation Status

| Category | Endpoints | Status |
|----------|-----------|--------|
| Repository Management | 45 | ✅ Complete |
| Pipeline & CI/CD | 38 | ✅ Complete |
| Issue Management | 25 | ✅ Complete |
| User & Team Management | 32 | ✅ Complete |
| Code Search & Snippets | 18 | ✅ Complete |
| Advanced Features | 20 | ✅ Complete |
| **Total** | **178** | **✅ 100%** |

## 🚀 Recent Updates

### Version 2.0.0 (December 2024)
- **Complete API Coverage**: All 178 endpoints implemented
- **Bug Fixes**: Critical nil pointer dereference issues resolved
- **OAuth Support**: Full OAuth 2.0 client credentials flow
- **New Features**: Snippets, GPG keys, code search, advanced pipeline management

## 🔧 Technical Details

### Swagger Specification
- **Format**: OpenAPI 3.0 (JSON)
- **Base URL**: `https://api.bitbucket.org/2.0/`
- **Authentication**: OAuth 2.0, Basic Auth, App Passwords
- **Rate Limiting**: 1000 requests/hour per user

### Provider Implementation
- **Framework**: Terraform Plugin SDK v2
- **Language**: Go 1.19+
- **Testing**: Unit tests + Acceptance tests
- **Documentation**: Auto-generated from schemas

---

**Note**: This swagger specification is maintained as the single source of truth for API implementation. Any changes to the provider should reference this specification to ensure consistency and completeness.
