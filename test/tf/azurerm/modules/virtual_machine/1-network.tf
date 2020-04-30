resource "azurerm_public_ip" "external" {
  name                = "${var.name}-pip"
  location            = var.location
  resource_group_name = var.resource_group_name
  allocation_method   = "Static"
  domain_name_label   = var.name
}

resource "azurerm_network_interface" "internal" {
  resource_group_name = var.resource_group_name
  location            = var.location
  name                = "${var.resource_group_name}-${random_string.resource_name.result}"

  ip_configuration {
    name                          = "private"
    subnet_id                     = var.subnet_id
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = azurerm_public_ip.external.id
  }

  tags = var.tags
}

resource "random_string" "resource_name" {
  length  = 16
  special = false
  upper   = false
  number  = false
}

