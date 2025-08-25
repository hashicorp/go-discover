# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "openstack" {
  version = "~> 1.27.0"
}

resource "openstack_compute_instance_v2" "consul-server" {
  count             = 3
  name              = "test-discover-${format("%01d", count.index + 1)}"
  image_name        = var.image
  flavor_name       = var.flavor
  availability_zone = var.az
  security_groups   = ["default"]

  metadata = {
    consul = "server"
  }

  network {
    uuid = var.network_uuid
  }
}

variable "image" {
}

variable "flavor" {
}

variable "az" {
}

variable "network_uuid" {
}

