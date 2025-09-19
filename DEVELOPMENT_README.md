# 🚀 Bitbucket Terraform Provider - Development Guide

## 🎉 **Version 2.0.0 - Complete API Coverage**

This provider now implements **100% of the Bitbucket API v3 specification** with 178 endpoints covering all Bitbucket Cloud functionality.

- **Website**: https://www.terraform.io
- **API Specification**: [`reference/swagger.v3.json`](reference/swagger.v3.json)
- **Coverage**: 178/178 endpoints (100% complete)
- **Status**: Production Ready

## 📋 Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19 (to build the provider plugin)
- [Git](https://git-scm.com/) for version control

## 🏗️ Building The Provider

### Quick Start
```bash
git clone https://github.com/gilesgamon/terraform-provider-bitbucket.git
cd terraform-provider-bitbucket
go build -o terraform-provider-bitbucket
```

### Development Setup
```bash
# Clone the repository
git clone https://github.com/gilesgamon/terraform-provider-bitbucket.git
cd terraform-provider-bitbucket

# Install dependencies
go mod download

# Build the provider
go build -o terraform-provider-bitbucket

# Run tests
go test ./...

# Format code
go fmt ./...

# Run linter
go vet ./...
```

## 🔧 Using the Provider

### Provider Configuration
```hcl
terraform {
  required_providers {
    bitbucket = {
      source  = "gilesgamon/terraform-provider-bitbucket"
      version = "2.0.0"
    }
  }
}

# Configure the Bitbucket Provider
provider "bitbucket" {
  # Option 1: Username/Password
  username = "your-username"
  password = "your-app-password"
  
  # Option 2: OAuth Client Credentials
  # oauth_client_id     = "your-client-id"
  # oauth_client_secret = "your-client-secret"
  
  # Option 3: OAuth Token
  # oauth_token = "your-oauth-token"
}
```

### Example Resources
```hcl
# Manage your repository
resource "bitbucket_repository" "infrastructure" {
  workspace = "myworkspace"
  name      = "terraform-code"
  project_key = "INFRA"
}

# Manage your project
resource "bitbucket_project" "infrastructure" {
  workspace = "myworkspace"
  name      = "Infrastructure"
  key       = "INFRA"
}

# Create a code snippet
resource "bitbucket_snippet" "example" {
  workspace = "myworkspace"
  title     = "Example Code"
  files = {
    "main.py" = "print('Hello, World!')"
  }
}
```

## 🛠️ Developing the Provider

### Prerequisites
- [Go](https://golang.org/doc/install) >= 1.19
- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Git](https://git-scm.com/)

### Development Workflow

#### 1. Code Quality Checks
```bash
# Format code
go fmt ./...

# Run linter
go vet ./...

# Check for issues
go mod tidy
```

#### 2. Testing
```bash
# Run unit tests
go test ./...

# Run acceptance tests (requires TF_ACC=1)
export TF_ACC=1
go test -v -timeout 120m ./...
```

#### 3. Building
```bash
# Build the provider
go build -o terraform-provider-bitbucket

# Test the build
./terraform-provider-bitbucket version
```

### 🐛 Recent Bug Fixes

#### Critical Fixes in v2.0.0
- **Nil Pointer Dereference**: Fixed critical crash in `bitbucket_repository` resource
- **Schema Validation**: Fixed `bitbucket_snippet` ID field type compliance
- **Type Safety**: Resolved compilation errors and type conflicts
- **Error Handling**: Improved error handling across all resources

### 📚 API Reference

The complete Bitbucket API v3 specification is available in [`reference/swagger.v3.json`](reference/swagger.v3.json). This provider implements all 178 endpoints from this specification.

### 🔍 Code Structure

```
bitbucket/
├── provider.go              # Provider configuration and schema
├── client.go                # HTTP client and API communication
├── resource_*.go            # Terraform resources (92 files)
├── data_*.go                # Terraform data sources (86 files)
└── docs/                    # Documentation for each resource/data source
```

### 🚀 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests: `go test ./...`
5. Format code: `go fmt ./...`
6. Submit a pull request
