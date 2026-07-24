# Copyright IBM Corp. 2017, 2026
# SPDX-License-Identifier: MPL-2.0

terraform {
  required_providers {
    azurerm = {
      source  = "registry.terraform.io/hashicorp/azurerm"
      version = "~> 4.81.0"
    }
    random = {
      source  = "registry.terraform.io/hashicorp/random"
      version = "~> 3.8.1"
    }
  }
}
provider "azurerm" {
  features {}
}

provider "random" {
}


variable "prefix" {
  default = "go-discover-azure-vmss"
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-dev"
  location = "West Europe"
}

module "network" {
  source         = "./modules/network"
  name           = "${var.prefix}-internalnw"
  resource_group = azurerm_resource_group.test.name
  location       = azurerm_resource_group.test.location
  address_space  = "10.0.0.0/16"
  subnet_cidr    = ["10.0.1.0/24"]
}

module "vmss" {
  source         = "./modules/vmss"
  name           = "${var.prefix}-01"
  resource_group = azurerm_resource_group.test.name
  location       = azurerm_resource_group.test.location
  subnet_id      = module.network.subnet_id
  size           = "Standard_DC1ds_v3"
}

