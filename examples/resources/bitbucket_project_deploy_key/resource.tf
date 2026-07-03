# Project deploy keys are inherited by every repository in the project,
# including repositories created outside Terraform.
resource "bitbucket_project_deploy_key" "shared" {
  workspace   = "my-workspace"
  project_key = "INFRA"
  key         = "ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAA..."
  label       = "shared-read-key"
}
