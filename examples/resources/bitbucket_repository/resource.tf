resource "bitbucket_project" "infra" {
  owner = "my-workspace"
  name  = "Infrastructure"
  key   = "INFRA"
}

resource "bitbucket_repository" "app" {
  owner       = "my-workspace"
  name        = "app"
  project_key = bitbucket_project.infra.key
  is_private  = true
}
