# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

//Configuring the provider
provider "scaleway" {
  version = "~> 1.15.0"
  region  = var.region
  zone    = var.zone
}

resource "scaleway_instance_server" "test" {
  count = 2
  image = var.image
  type  = "DEV1-S"
  tags  = ["consul-server"]
  name  = "test-server"
  zone  = var.zone
}

resource "scaleway_instance_server" "dummy" {
  image = var.image
  name  = "dummy_scaleway_instance"
  type  = "DEV1-S"
  tags  = ["Dummy"]
  zone  = var.zone
}

