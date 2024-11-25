# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "alicloud" {
  version = "1.80.1"
}

resource "alicloud_instance" "test" {
  count           = 2
  image_id        = var.image_id
  instance_type   = var.instance_type
  security_groups = [alicloud_security_group.default.id]
  tags            = var.tag
  vswitch_id      = alicloud_vswitch.vswitch.id
}

resource "alicloud_security_group" "default" {
  name        = "default"
  description = "default"
  vpc_id      = alicloud_vpc.vpc.id
}

resource "alicloud_vpc" "vpc" {
  name       = "go_discover_test"
  cidr_block = "10.1.0.0/24"
}

resource "alicloud_vswitch" "vswitch" {
  vpc_id            = alicloud_vpc.vpc.id
  cidr_block        = alicloud_vpc.vpc.cidr_block
  availability_zone = "us-west-1a"
}