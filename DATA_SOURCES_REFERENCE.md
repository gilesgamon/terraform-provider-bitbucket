# Bitbucket Provider Data Sources Reference

## 🚀 **Quick Reference Card**

### **Core Repository & Git Operations** ✅

| Data Source | Purpose | Key Parameters | Key Outputs |
|-------------|---------|----------------|-------------|
| `bitbucket_repository` | Repository information | `workspace`, `repo_slug` | `name`, `description`, `is_private` |
| `bitbucket_branch` | Branch details | `workspace`, `repo_slug`, `branch_name` | `target.hash`, `name`, `type` |
| `bitbucket_tag` | Tag information | `workspace`, `repo_slug`, `tag_name` | `target.hash`, `name`, `type` |
| `bitbucket_commit` | Commit details | `workspace`, `repo_slug`, `commit_sha` | `hash`, `author`, `message`, `date` |
| `bitbucket_pull_request` | PR information | `workspace`, `repo_slug`, `pull_request_id` | `title`, `state`, `author`, `source` |
| `bitbucket_pipeline` | Pipeline details | `workspace`, `repo_slug`, `pipeline_number` | `state`, `trigger`, `target`, `created_on` |

### **Advanced Git Operations** ✅

| Data Source | Purpose | Key Parameters | Key Outputs |
|-------------|---------|----------------|-------------|
| `bitbucket_commit_comments` | Commit comments | `workspace`, `repo_slug`, `commit_sha` | `comments[]` |
| `bitbucket_commit_statuses` | Build statuses | `workspace`, `repo_slug`, `commit_sha` | `statuses[]` |
| `bitbucket_commit_properties` | Commit metadata | `workspace`, `repo_slug`, `commit_sha`, `app_key` | `properties[]` |
| `bitbucket_commit_reports` | Quality reports | `workspace`, `repo_slug`, `commit_sha` | `reports[]` |
| `bitbucket_commit_pullrequests` | PRs with commit | `workspace`, `repo_slug`, `commit_sha` | `pull_requests[]` |
| `bitbucket_commit_approvals` | Commit approvals | `workspace`, `repo_slug`, `commit_sha` | `approvals[]` |
| `bitbucket_commits` | Commit list | `workspace`, `repo_slug` | `commits[]` |
| `bitbucket_commit_diff` | Commit changes | `workspace`, `repo_slug`, `commit_sha` | `diff[]` |
| `bitbucket_commit_diffstat` | Change statistics | `workspace`, `repo_slug`, `commit_sha` | `diffstat[]`, `summary` |

### **Issue Management** 🔄

| Data Source | Purpose | Key Parameters | Key Outputs |
|-------------|---------|----------------|-------------|
| `bitbucket_issues` | Issue list | `workspace`, `repo_slug`, `state`, `kind`, `priority` | `issues[]` |

---

## 📋 **Common Parameters**

### **Required Parameters**
- `workspace`: Your Bitbucket workspace/team name
- `repo_slug`: Repository name/slug

### **Optional Filters**
- `state`: Filter by state (open, resolved, closed, declined, merged)
- `kind`: Filter by kind (bug, enhancement, proposal, task)
- `priority`: Filter by priority (trivial, minor, major, critical, blocker)
- `assignee`: Filter by assignee username
- `reporter`: Filter by reporter username
- `milestone`: Filter by milestone name
- `component`: Filter by component name
- `version`: Filter by version name
- `q`: Search query string
- `sort`: Sort field (created_on, updated_on, priority, kind, state)

---

## 🔄 **Usage Patterns**

### **Pipeline Integration**
```hcl
# Get latest commit for pipeline
data "bitbucket_branch" "main" {
  workspace = "myworkspace"
  repo_slug = "my-app"
  branch_name = "main"
}

# Use in pipeline
resource "aws_codepipeline" "app" {
  # ... other config
  source {
    revision = data.bitbucket_branch.main.target.hash
  }
}
```

### **Issue Monitoring**
```hcl
# Monitor critical issues
data "bitbucket_issues" "critical" {
  workspace = "myworkspace"
  repo_slug = "my-app"
  state = "open"
  kind = "bug"
  priority = "critical"
}

output "critical_issue_count" {
  value = length(data.bitbucket_issues.critical.issues)
}
```

### **Commit Analysis**
```hcl
# Analyze commit impact
data "bitbucket_commit_diffstat" "changes" {
  workspace = "myworkspace"
  repo_slug = "my-app"
  commit_sha = "abc123..."
}

output "impact_metrics" {
  value = {
    files_changed = data.bitbucket_commit_diffstat.changes.summary["total_files"]
    lines_added = data.bitbucket_commit_diffstat.changes.summary["total_added"]
    lines_removed = data.bitbucket_commit_diffstat.changes.summary["total_removed"]
  }
}
```

---

## 🎯 **Best Practices**

1. **Use Workspace Variables**: Store workspace names in variables
2. **Filter Early**: Use filters to reduce API calls and data transfer
3. **Cache Results**: Use `depends_on` to control when data sources refresh
4. **Error Handling**: Always check for required outputs before using them
5. **Documentation**: Document complex queries and their business logic

---

## 📚 **Full Documentation**

For complete documentation, examples, and advanced usage patterns, see:
- [Main README](README.md)
- [Implementation Progress](IMPLEMENTATION_PROGRESS.md)
- [Terraform Registry](https://registry.terraform.io/providers/DrFaust92/bitbucket/latest/docs)
