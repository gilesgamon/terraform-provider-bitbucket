terraform {
  required_providers {
    bitbucket = {
      source  = "gilesgamon/terraform-provider-bitbucket"
      version = "~> 0.1"
    }
  }
}

# Authenticate with a username and app password.
# Credentials can also be supplied via the BITBUCKET_USERNAME and
# BITBUCKET_PASSWORD environment variables.
provider "bitbucket" {
  username = "my-user"
  password = "my-app-password"
}
