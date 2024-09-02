# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

module "network" {
  source        = "./modules/network"
  address_space = var.address_space
  subnet_cidr   = var.subnet_cidr
}

data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "tagged" {
  count                       = 2
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t2.micro"
  subnet_id                   = module.network.subnet_id
  associate_public_ip_address = true

  tags = {
    "consul" = "server"
  }
}

// We keep an extra untagged resource to test we get back
// 2/3 of our instances
resource "aws_instance" "not-tagged" {
  ami                         = data.aws_ami.ubuntu.id
  instance_type               = "t2.micro"
  subnet_id                   = module.network.subnet_id
  associate_public_ip_address = true
}

