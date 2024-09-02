# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "azurerm" {
  version = "~> 2.7.0"
  features {}
}

provider "random" {
  version = "~> 2.2.1"
}


variable "prefix" {
  default = "go-discover-azurerm"
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-dev"
  location = "West Europe"
}

module "network" {
  source              = "./modules/network"
  name                = "${var.prefix}-internalnw"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  address_space       = "10.0.0.0/16"
  subnet_cidr         = "10.0.1.0/24"
}

module "vm01" {
  source              = "./modules/virtual_machine"
  name                = "${var.prefix}-01"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = module.network.subnet_id

  tags = {
    "consul" = "server"
  }
}

module "vm02" {
  source              = "./modules/virtual_machine"
  name                = "${var.prefix}-02"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = module.network.subnet_id

  tags = {
    "consul" = "server"
  }
}

// We intentionally don't tag the last machine to ensure we only discover the
// first two
module "vm03" {
  source              = "./modules/virtual_machine"
  name                = "${var.prefix}-03"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location
  subnet_id           = module.network.subnet_id
}

