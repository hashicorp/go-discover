variable "prefix" {
  default = "go-discover"
}

resource "azurerm_resource_group" "test" {
  name     = "${var.prefix}-dev"
  location = "West Europe"
}

module "network" {
  source              = "./modules/network"
  name                = "${var.prefix}-internalnw"
  resource_group_name = "${azurerm_resource_group.test.name}"
  location            = "${azurerm_resource_group.test.location}"
  address_space       = "10.0.0.0/16"
  subnet_cidr         = "10.0.1.0/24"
}

module "vm01" {
  source               = "./modules/virtual_machine"
  name                 = "${var.prefix}-01"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  subnet_id            = "${module.network.subnet_id}"

  tags {
    "consul" = "server"
  }
}

module "vm02" {
  source               = "./modules/virtual_machine"
  name                 = "${var.prefix}-02"
  resource_group_name  = "${azurerm_resource_group.test.name}"
  location             = "${azurerm_resource_group.test.location}"
  subnet_id            = "${module.network.subnet_id}"
}

output "public_ips" {
  value = ["${module.vm01.public_ip}","${module.vm02.public_ip}"]
}

output "private_ips" {
  value = ["${module.vm01.private_ip}", "${module.vm02.private_ip}"]
}

output "tagged_ips" {
  value = ["${module.vm01.private_ip}"]
}
