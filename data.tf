provider "vault" {
  address= var.VAULT_ADDR
	token = var.ROOT_TOKEN
}

data "vault_generic_secret" "aws_secrets"{
  path = "secret/aws"
}

data "aws_availability_zones" "available" {}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-focal-20.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }
  
  owners = ["099720109477"] # Canonical
}