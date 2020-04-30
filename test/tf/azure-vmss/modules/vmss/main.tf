provider "azurerm" {
  version = "~> 2.7.0"
  features {}
}

provider "random" {
  version = "~> 2.2.1"
}

resource "azurerm_public_ip" "test" {
  name                = "${var.name}-pip"
  location            = var.location
  resource_group_name = var.resource_group
  allocation_method   = "Static"
  domain_name_label   = var.resource_group

  tags = var.tags
}

resource "azurerm_lb" "test" {
  name                = "${var.name}-lb"
  location            = var.location
  resource_group_name = var.resource_group

  frontend_ip_configuration {
    name                 = "PublicIPAddress"
    public_ip_address_id = azurerm_public_ip.test.id
  }
}

resource "azurerm_lb_backend_address_pool" "bpepool" {
  resource_group_name = var.resource_group
  loadbalancer_id     = azurerm_lb.test.id
  name                = "${var.resource_group}-${random_string.resource_name.result}"
}

resource "azurerm_virtual_machine_scale_set" "test" {
  name                = "${var.name}-scale-set"
  location            = var.location
  resource_group_name = var.resource_group
  upgrade_policy_mode = "Manual"
  overprovision       = false

  sku {
    name     = var.size
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
    computer_name_prefix = var.name
    admin_username       = var.username
    admin_password       = random_string.password.result
  }

  network_profile {
    name    = "${var.name}-np"
    primary = true

    ip_configuration {
      name                                   = "${var.name}-ipc"
      primary                                = true
      subnet_id                              = var.subnet_id
      load_balancer_backend_address_pool_ids = [azurerm_lb_backend_address_pool.bpepool.id]
    }
  }

  tags = var.tags
}

resource "random_string" "password" {
  length = 16
}

resource "random_string" "resource_name" {
  length  = 16
  special = false
  upper   = false
  number  = false
}

