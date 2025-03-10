terraform{
  required_version = ">= 1.0.3"
  required_providers {
    aws = {
      source = "hashicorp/aws"
      version = "5.82.2"
    }
    local = {
      source = "hashicorp/local"
      version = "2.5.2"
    }
  }
}

provider "aws" {
  region = var.AWS_REGION
  access_key = data.vault_generic_secret.aws_secrets.data["access_key"]
  secret_key = data.vault_generic_secret.aws_secrets.data["secret_key"]
}