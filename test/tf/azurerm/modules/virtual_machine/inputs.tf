# Copyright IBM Corp. 2017, 2026
# SPDX-License-Identifier: MPL-2.0

variable "name" {
}

variable "resource_group_name" {
}

variable "location" {
}

variable "subnet_id" {
}

variable "size" {
  default = "Standard_DC1ds_v3"
}

variable "username" {
  default = "ubuntu"
}

variable "tags" {
  type    = map(string)
  default = {}
}

