# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

variable "tag" {
  type = map(string)

  default = {
    "consul" = "server.test"
  }
}

variable "instance_type" {
  default = "ecs.n4.small"
}

variable "image_id" {
  default = "centos_7_04_64_20G_alibase_201701015.vhd"
}

