
locals {
  instance_type = var.instance_type
  tag_name = var.role
}

resource "aws_instance" "instance" {
  count = var.instance_count
  ami = var.ec2_ami
  instance_type = local.instance_type
  availability_zone = var.availability_zone
  associate_public_ip_address = true
  vpc_security_group_ids = var.sg_ids
  key_name = var.key_name
  tags = {
    Name = "${local.tag_name}-instance-${count.index}"
    Environment = var.environment
    Role = var.role
  }
}