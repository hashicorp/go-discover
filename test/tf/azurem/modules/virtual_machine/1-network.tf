resource "azurerm_public_ip" "external" {
  name                         = "${var.name}-pip"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group_name}"
  public_ip_address_allocation = "static"
  domain_name_label            = "${var.name}"
}

resource "azurerm_network_interface" "internal" {
  name                = "${var.name}-nic"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group_name}"

  ip_configuration {
    name                          = "private"
    subnet_id                     = "${var.subnet_id}"
    private_ip_address_allocation = "dynamic"
    public_ip_address_id          = "${azurerm_public_ip.external.id}"
  }

  tags = "${var.tags}"
}
