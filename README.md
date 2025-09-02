# Bitbucket Terraform Provider

[![Go Report Card](https://goreportcard.com/badge/github.com/DrFaust92/bitbucket)](https://goreportcard.com/report/github.com/DrFaust92/bitbucket)
[![GoDoc](https://godoc.org/github.com/DrFaust92/bitbucket?status.svg)](https://godoc.org/github.com/DrFaust92/bitbucket)
[![License](https://img.shields.io/github/license/DrFaust92/bitbucket.svg)](https://github.com/DrFaust92/bitbucket/blob/master/LICENSE)

The Bitbucket Terraform Provider enables you to manage your Bitbucket Cloud resources using Terraform. This provider offers comprehensive coverage of Bitbucket's API, including repositories, pipelines, issues, and advanced Git operations.

## 🚀 Features

### ✅ **Core Repository & Git Operations (100% Complete)**
- **Repository Management**: Create, update, and manage repositories
- **Branch Operations**: Manage branches, restrictions, and branching models
- **Tag Management**: Handle Git tags and releases
- **Commit Operations**: Access commit details, properties, reports, and approvals
- **Pull Request Management**: Comprehensive PR lifecycle management
- **Pipeline Integration**: CI/CD pipeline management and monitoring

### 🔄 **Advanced Repository Features (In Progress)**
- **Issue Tracking**: Complete issue management system
- **Repository Settings**: Configuration and permissions
- **Webhook Management**: Event-driven integrations
- **Deploy Keys**: SSH key management for deployments

### 🔄 **Pipeline & CI/CD (In Progress)**
- **Build Management**: Pipeline runs, steps, and logs
- **Environment Configuration**: Deployment environments and variables
- **Test Reports**: Automated testing and reporting

### 🔄 **Workspace & Project Management (Planned)**
- **Team Management**: Workspace and project organization
- **Permission Management**: Access control and security
- **Variable Management**: Environment and repository variables

## 📋 Requirements

- [Terraform](https://www.terraform.io/downloads.html) >= 1.0
- [Go](https://golang.org/doc/install) >= 1.19 (for building from source)

## 🔧 Installation

### Using Terraform Registry (Recommended)

```hcl
terraform {
  required_providers {
    bitbucket = {
      source  = "DrFaust92/bitbucket"
      version = "~> 2.45.1"
    }
  }
}
```

### Building from Source

```bash
git clone https://github.com/DrFaust92/bitbucket.git
cd bitbucket
go build -o terraform-provider-bitbucket .
```

## ⚙️ Provider Configuration

### Basic Authentication

```hcl
provider "bitbucket" {
  username = "your-username"
  password = "your-password-or-app-password"
}
```

### App Password (Recommended)

```hcl
provider "bitbucket" {
  username = "your-username"
  password = "your-app-password"
}
```

### OAuth2 (Advanced)

```hcl
provider "bitbucket" {
  client_id     = "your-client-id"
  client_secret = "your-client-secret"
  token         = "your-oauth-token"
}
```

## 📚 Data Sources

### Repository & Git Operations

#### Get Repository Information
```hcl
data "bitbucket_repository" "main" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
}

output "repo_name" {
  value = data.bitbucket_repository.main.name
}
```

#### Get Commit Details
```hcl
data "bitbucket_commit" "latest" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
}

output "commit_author" {
  value = data.bitbucket_commit.latest.author
}
```

#### Get Branch Information
```hcl
data "bitbucket_branch" "main" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  branch_name = "main"
}

output "latest_commit" {
  value = data.bitbucket_branch.main.target.hash
}
```

#### Get Tag Information
```hcl
data "bitbucket_tag" "release" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  tag_name = "v1.0.0"
}

output "tag_commit" {
  value = data.bitbucket_tag.release.target.hash
}
```

#### Get Pull Request Details
```hcl
data "bitbucket_pull_request" "feature" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  pull_request_id = "123"
}

output "pr_title" {
  value = data.bitbucket_pull_request.feature.title
}
```

#### Get Pipeline Information
```hcl
data "bitbucket_pipeline" "build" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  pipeline_number = "456"
}

output "pipeline_state" {
  value = data.bitbucket_pipeline.build.state.name
}
```

### Advanced Git Operations

#### Get Commit Comments
```hcl
data "bitbucket_commit_comments" "feedback" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
}

output "comment_count" {
  value = length(data.bitbucket_commit_comments.feedback.comments)
}
```

#### Get Commit Statuses
```hcl
data "bitbucket_commit_statuses" "checks" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
}

output "status_count" {
  value = length(data.bitbucket_commit_statuses.checks.statuses)
}
```

#### Get Commit Properties
```hcl
data "bitbucket_commit_properties" "metadata" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
  app_key = "my-app"
}

output "property_count" {
  value = length(data.bitbucket_commit_properties.metadata.properties)
}
```

#### Get Commit Reports
```hcl
data "bitbucket_commit_reports" "quality" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
}

output "report_count" {
  value = length(data.bitbucket_commit_reports.quality.reports)
}
```

#### Get Commits List
```hcl
data "bitbucket_commits" "recent" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  branch = "main"
}

output "recent_commits" {
  value = [for commit in data.bitbucket_commits.recent.commits : commit.hash]
}
```

#### Get Commit Diff
```hcl
data "bitbucket_commit_diff" "changes" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
  context = 3
}

output "files_changed" {
  value = length(data.bitbucket_commit_diff.changes.diff)
}
```

#### Get Commit Diff Statistics
```hcl
data "bitbucket_commit_diffstat" "stats" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  commit_sha = "abc123..."
}

output "lines_added" {
  value = data.bitbucket_commit_diffstat.stats.summary["total_added"]
}
```

### Issue Management

#### Get Issues List
```hcl
data "bitbucket_issues" "bugs" {
  workspace = "myworkspace"
  repo_slug = "my-repo"
  state = "open"
  kind = "bug"
  priority = "major"
}

output "open_bugs" {
  value = [for issue in data.bitbucket_issues.bugs.issues : issue.title]
}
```

## 🏗️ Resources

### Repository Management

```hcl
resource "bitbucket_repository" "infrastructure" {
  workspace = "myworkspace"
  name     = "terraform-infrastructure"
  project_key = "INFRA"
  
  is_private = true
  fork_policy = "allow_forks"
  
  description = "Infrastructure as Code repository"
}
```

### Project Management

```hcl
resource "bitbucket_project" "infrastructure" {
  workspace = "myworkspace"
  name     = "Infrastructure"
  key      = "INFRA"
  
  description = "Infrastructure and DevOps projects"
  is_private = true
}
```

### Pipeline Configuration

```hcl
resource "bitbucket_pipeline_variable" "environment" {
  workspace = "myworkspace"
  repository = "my-repo"
  key = "ENVIRONMENT"
  value = "production"
  secured = false
}
```

## 🔄 Usage Examples

### Complete CI/CD Pipeline Setup

```hcl
# Get repository information
data "bitbucket_repository" "app" {
  workspace = "myworkspace"
  repo_slug = "my-application"
}

# Get latest commit from main branch
data "bitbucket_branch" "main" {
  workspace = "myworkspace"
  repo_slug = "my-application"
  branch_name = "main"
}

# Get pipeline information
data "bitbucket_pipeline" "latest" {
  workspace = "myworkspace"
  repo_slug = "my-application"
  pipeline_number = "latest"
}

# Output pipeline metadata
output "pipeline_info" {
  value = {
    repository = data.bitbucket_repository.app.name
    latest_commit = data.bitbucket_branch.main.target.hash
    pipeline_state = data.bitbucket_pipeline.latest.state.name
    pipeline_created = data.bitbucket_pipeline.latest.created_on
  }
}
```

### Issue Tracking Integration

```hcl
# Get all open issues
data "bitbucket_issues" "open_issues" {
  workspace = "myworkspace"
  repo_slug = "my-application"
  state = "open"
}

# Get high-priority bugs
data "bitbucket_issues" "critical_bugs" {
  workspace = "myworkspace"
  repo_slug = "my-application"
  state = "open"
  kind = "bug"
  priority = "critical"
}

output "issue_summary" {
  value = {
    total_open = length(data.bitbucket_issues.open_issues.issues)
    critical_bugs = length(data.bitbucket_issues.critical_bugs.issues)
  }
}
```

## 🧪 Testing

### Run Unit Tests
```bash
go test ./...
```

### Run Acceptance Tests
```bash
export TF_ACC=1
go test -v -timeout 120m ./...
```

## 📖 Documentation

For detailed documentation on each data source and resource, see the [Terraform Registry](https://registry.terraform.io/providers/DrFaust92/bitbucket/latest/docs).

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

### Development Setup

1. Fork the repository
2. Clone your fork
3. Install dependencies: `go mod download`
4. Make your changes
5. Run tests: `go test ./...`
6. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

- **Issues**: [GitHub Issues](https://github.com/DrFaust92/bitbucket/issues)
- **Documentation**: [Terraform Registry](https://registry.terraform.io/providers/DrFaust92/bitbucket/latest/docs)
- **Discussions**: [GitHub Discussions](https://github.com/DrFaust92/bitbucket/discussions)

## 🔗 Related Links

- [Terraform Documentation](https://www.terraform.io/docs)
- [Bitbucket API Documentation](https://developer.atlassian.com/cloud/bitbucket/rest/)
- [Terraform Provider Development](https://www.terraform.io/docs/extend/index.html)

---

**Made with ❤️ by the Terraform community**
