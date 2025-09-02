# Enhanced Terraform Provider for Bitbucket

This enhanced version of the `terraform-provider-bitbucket` includes additional data sources that enable advanced pipeline integration and Git workflow automation.

## New Enhanced Data Sources

### 1. `bitbucket_tag` - Tag Information

Retrieves detailed information about a specific tag in a repository.

**Usage:**
```hcl
data "bitbucket_tag" "release_tag" {
  workspace = "your-workspace"
  repo_slug = "your-repository"
  tag_name  = "v1.0.0"
}
```

**Available Attributes:**
- `uuid` - Tag UUID
- `name` - Tag name
- `target_hash` - Commit SHA that the tag points to
- `target_date` - Date when the tag was created
- `message` - Tag message/description

**Use Cases:**
- Getting the commit SHA for a specific release tag
- Automating deployments based on tag information
- Pipeline integration for release management

### 2. `bitbucket_commit` - Commit Information

Retrieves detailed information about a specific commit or branch.

**Usage:**
```hcl
data "bitbucket_commit" "latest_commit" {
  workspace  = "your-workspace"
  repo_slug  = "your-repository"
  commit_sha = "main"  # Can be branch name or commit SHA
}
```

**Available Attributes:**
- `hash` - Commit SHA
- `message` - Commit message
- `date` - Commit date
- `author` - Author information (username, display_name, uuid)
- `parents` - Parent commit information

**Use Cases:**
- Getting the latest commit SHA from a branch
- Pipeline integration for commit-based deployments
- Audit and compliance tracking

### 3. `bitbucket_branch` - Branch Information

Retrieves detailed information about a specific branch.

**Usage:**
```hcl
data "bitbucket_branch" "main_branch" {
  workspace   = "your-workspace"
  repo_slug   = "your-repository"
  branch_name = "main"
}
```

**Available Attributes:**
- `name` - Branch name
- `target_hash` - Latest commit SHA on the branch
- `target_date` - Date of the latest commit
- `target_message` - Message of the latest commit

**Use Cases:**
- Getting the latest commit SHA from a branch
- Branch-based pipeline triggers
- Development workflow automation

### 4. `bitbucket_pull_request` - Pull Request Information

Retrieves detailed information about a specific pull request.

**Usage:**
```hcl
data "bitbucket_pull_request" "feature_pr" {
  workspace      = "your-workspace"
  repo_slug      = "your-repository"
  pull_request_id = "123"
}
```

**Available Attributes:**
- `id` - Pull request ID
- `title` - Pull request title
- `description` - Pull request description
- `state` - Current state (OPEN, MERGED, DECLINED)
- `author` - Author information
- `source` - Source branch and commit information
- `destination` - Destination branch and commit information
- `created_date` - Creation date
- `updated_date` - Last update date
- `merge_commit` - Merge commit information (if merged)

**Use Cases:**
- Pull request-based pipeline triggers
- Code review automation
- Merge workflow management

### 5. `bitbucket_pipeline` - Pipeline Information

Retrieves detailed information about a specific pipeline run.

**Usage:**
```hcl
data "bitbucket_pipeline" "build_pipeline" {
  workspace       = "your-workspace"
  repo_slug       = "your-repository"
  pipeline_number = "456"
}
```

**Available Attributes:**
- `id` - Pipeline ID
- `build_number` - Build number
- `state` - Pipeline state
- `created_on` - Creation timestamp
- `completed_on` - Completion timestamp
- `trigger` - Trigger information and user
- `target` - Target branch/commit information

**Use Cases:**
- Pipeline status monitoring
- Build artifact tracking
- Deployment automation

## Pipeline Integration Examples

### Example 1: Tag-Based Deployment

```hcl
# Get information about a specific release tag
data "bitbucket_tag" "release" {
  workspace = "my-company"
  repo_slug = "my-app"
  tag_name  = "v2.1.0"
}

# Use the tag's commit SHA in your pipeline
resource "aws_codepipeline" "deploy" {
  # ... other configuration ...
  
  stage {
    name = "Source"
    action {
      name     = "Source"
      category = "Source"
      owner    = "AWS"
      provider = "CodeCommit"
      version  = "1"
      
      configuration = {
        RepositoryName = "my-app-repo"
        BranchName     = "main"
        CommitId       = data.bitbucket_tag.release.target_hash
      }
    }
  }
}
```

### Example 2: Branch-Based Pipeline

```hcl
# Get the latest commit from a development branch
data "bitbucket_branch" "develop" {
  workspace   = "my-company"
  repo_slug   = "my-app"
  branch_name = "develop"
}

# Use the branch's latest commit in your pipeline
resource "aws_codepipeline" "dev_deploy" {
  # ... other configuration ...
  
  stage {
    name = "Source"
    action {
      name     = "Source"
      category = "Source"
      owner    = "AWS"
      provider = "CodeCommit"
      version  = "1"
      
      configuration = {
        RepositoryName = "my-app-repo"
        BranchName     = "develop"
        CommitId       = data.bitbucket_branch.develop.target_hash
      }
    }
  }
}
```

### Example 3: Pull Request Integration

```hcl
# Get information about a specific pull request
data "bitbucket_pull_request" "feature" {
  workspace      = "my-company"
  repo_slug      = "my-app"
  pull_request_id = "123"
}

# Use PR information in your pipeline
resource "aws_codepipeline" "pr_deploy" {
  # ... other configuration ...
  
  stage {
    name = "Source"
    action {
      name     = "Source"
      category = "Source"
      owner    = "AWS"
      provider = "CodeCommit"
      version  = "1"
      
      configuration = {
        RepositoryName = "my-app-repo"
        BranchName     = data.bitbucket_pull_request.feature.source[0].branch
        CommitId       = data.bitbucket_pull_request.feature.source[0].commit
      }
    }
  }
}
```

## Authentication

The provider supports multiple authentication methods:

### Username/Password (App Password)
```hcl
provider "bitbucket" {
  username = "your-username"
  password = "your-app-password"
}
```

### OAuth Token
```hcl
provider "bitbucket" {
  oauth_token = "your-oauth-token"
}
```

### OAuth Client Credentials
```hcl
provider "bitbucket" {
  oauth_client_id     = "your-client-id"
  oauth_client_secret = "your-client-secret"
}
```

## Environment Variables

You can also use environment variables for authentication:

```bash
export BITBUCKET_USERNAME="your-username"
export BITBUCKET_PASSWORD="your-app-password"
# or
export BITBUCKET_OAUTH_TOKEN="your-oauth-token"
# or
export BITBUCKET_OAUTH_CLIENT_ID="your-client-id"
export BITBUCKET_OAUTH_CLIENT_SECRET="your-client-secret"
```

## Building the Provider

To build the enhanced provider:

```bash
cd /path/to/terraform-provider-bitbucket
go build -o terraform-provider-bitbucket .
```

## Testing

Use the provided `test_enhanced_provider.tf` file to test the enhanced data sources:

```bash
terraform init
terraform plan
```

## Use Cases

These enhanced data sources enable:

1. **Automated Deployments**: Use tag or branch information to automatically deploy specific versions
2. **Pipeline Integration**: Integrate Bitbucket information into AWS CodePipeline, GitHub Actions, or other CI/CD tools
3. **Release Management**: Automate release processes based on tag information
4. **Branch Protection**: Use branch information for automated testing and validation
5. **Pull Request Automation**: Automate workflows based on PR status and information
6. **Audit and Compliance**: Track commit history and changes for compliance purposes

## Contributing

When adding new data sources or resources:

1. Follow the existing code patterns
2. Use the HTTP client approach for API calls
3. Implement proper error handling
4. Add comprehensive documentation
5. Include test configurations

## Support

For issues or questions:
1. Check the existing data source implementations for reference
2. Review the Bitbucket REST API documentation
3. Test with the provided example configurations
