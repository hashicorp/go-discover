output "public_ip" {
  value = azurerm_public_ip.test.ip_address
}

output "vm_scale_set" {
  value = azurerm_virtual_machine_scale_set.test.name
}

