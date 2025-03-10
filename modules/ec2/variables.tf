variable "instance_type" {}
variable "role" {}
variable "instance_count" {}
variable "ec2_ami"{}
variable "availability_zone"{}
variable "key_name"{}
variable "environment"{}
variable "sg_ids" {
  type = list(string)
}