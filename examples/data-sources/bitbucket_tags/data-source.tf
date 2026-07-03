data "bitbucket_tags" "example" {
  workspace = "my-workspace"
  repo_slug = "app"
}

output "tag_names" {
  value = [for t in data.bitbucket_tags.example.tags : t.name]
}
