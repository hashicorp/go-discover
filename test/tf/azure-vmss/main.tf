provider "azurerm" {
  version = "~> 2.7.0"
  features {}
}

provider "random" {
  version = "~> 2.2.1"
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
  subnet_cidr    = "10.0.1.0/24"
}

module "vmss" {
  source         = "./modules/vmss"
  name           = "${var.prefix}-01"
  resource_group = azurerm_resource_group.test.name
  location       = azurerm_resource_group.test.location
  subnet_id      = module.network.subnet_id
}

