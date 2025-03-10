locals {
  key_name ="ansible_key"
}

module "sec_group" {
source = "./modules/security-group"
}

module "ssh_key" {
source = "./modules/ssh-keys"
private_key_filepath = var.private_key_filepath
key_name = local.key_name
}

module "web" {
source = "./modules/ec2"
instance_type = "t2.micro"
role = "web"
instance_count = 2
ec2_ami = data.aws_ami.ubuntu.id
availability_zone = data.aws_availability_zones.available.names[0]
key_name = local.key_name
environment = var.environment
sg_ids = [module.sec_group.security_group_id]
}

module "database" {
source = "./modules/ec2"
instance_type = "t2.micro"
role = "database"
instance_count = 2
ec2_ami = data.aws_ami.ubuntu.id
availability_zone = data.aws_availability_zones.available.names[0]
key_name = local.key_name
environment = var.environment
sg_ids = [module.sec_group.security_group_id]
}

module "application" {
source = "./modules/ec2"
instance_type = "t2.micro"
role = "application"
instance_count = 2
ec2_ami = data.aws_ami.ubuntu.id
availability_zone = data.aws_availability_zones.available.names[0]
key_name = local.key_name
environment = var.environment
sg_ids = [module.sec_group.security_group_id]
}
