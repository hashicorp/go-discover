# Copyright IBM Corp. 2017, 2025
# SPDX-License-Identifier: MPL-2.0

output "public_ip" {
  value = azurerm_public_ip.external.ip_address
}

output "private_ip" {
  value = azurerm_network_interface.internal.private_ip_address
}

