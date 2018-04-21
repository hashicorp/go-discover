resource "azurerm_public_ip" "test" {
  name                         = "${var.name}-pip"
  location                     = "${var.location}"
  resource_group_name          = "${var.resource_group}"
  public_ip_address_allocation = "static"
  domain_name_label            = "${var.resource_group}"

  tags = "${var.tags}"
}

resource "azurerm_lb" "test" {
  name                = "${var.name}-lb"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group}"

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = "${azurerm_public_ip.test.id}"
  }
}

resource "azurerm_lb_backend_address_pool" "bpepool" {
  resource_group_name = "${var.resource_group}"
  loadbalancer_id     = "${azurerm_lb.test.id}"
  name                = "BackEndAddressPool"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "${var.name}-vmss"
  location            = "${var.location}"
  resource_group_name = "${var.resource_group}"
  upgrade_policy_mode = "Manual"

  sku {
    name     = "${var.size}"
    tier     = "Standard"
    capacity = 3
  }

  storage_profile_image_reference {
    publisher = "Canonical"
    offer     = "UbuntuServer"
    sku       = "16.04-LTS"
    version   = "latest"
  }

  storage_profile_os_disk {
    name              = ""
    caching           = "ReadWrite"
    create_option     = "FromImage"
    managed_disk_type = "Standard_LRS"
  }

  os_profile {
    computer_name_prefix = "${var.name}"
    admin_username       = "${var.username}"
    admin_password       = ""
  }

  network_profile {
    name    = "${var.name}-np"
    primary = true

    ip_configuration {
      name                                   = "${var.name}-ipc"
      subnet_id                              = "${var.subnet_id}"
      load_balancer_backend_address_pool_ids = ["${azurerm_lb_backend_address_pool.bpepool.id}"]
    }
  }

  tags = "${var.tags}"
}
