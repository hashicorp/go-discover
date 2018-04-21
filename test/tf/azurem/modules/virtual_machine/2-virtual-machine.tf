resource "azurerm_virtual_machine" "main" {
  name                  = "${var.name}-vm"
  location              = "${var.location}"
  resource_group_name   = "${var.resource_group_name}"
  network_interface_ids = ["${azurerm_network_interface.internal.id}"]
  vm_size               = "${var.size}"

  storage_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "14.04.5-LTS"
    version   = "latest"
  }

  storage_os_disk {
    name              = "${var.name}-osdisk"
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }
}
