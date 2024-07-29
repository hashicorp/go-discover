# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "tencentcloud" {
  version = ">= 1.32.1"
}

resource "tencentcloud_instance" "test" {
  count                      = 2
  instance_name              = "test"
  availability_zone          = var.availability_zone
  image_id                   = var.image_id
  instance_type              = var.instance_type
  system_disk_type           = "CLOUD_PREMIUM"
  system_disk_size           = 50
  internet_max_bandwidth_out = 1
  allocate_public_ip         = true
  tags                       = var.tag
}

