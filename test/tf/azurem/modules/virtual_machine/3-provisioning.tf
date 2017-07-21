resource "null_resource" "provisioning" {
  connection {
    type        = "ssh"
    user        = "${var.username}"
    host        = "${azurerm_public_ip.external.ip_address}"
    timeout     = "10m"
    private_key = "${file(var.private_ssh_key_path)}"
  }

  provisioner "file" {
    source      = "discover"
    destination = "/home/${var.username}/discover"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /home/${var.username}/discover"
    ]
  }

  depends_on = ["azurerm_virtual_machine.main"]
}
