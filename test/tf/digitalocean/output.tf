output "public_ips" {
  value = ["${digitalocean_droplet.test-01.ipv4_address }","${digitalocean_droplet.test-02.ipv4_address }"]
}

output "private_ips" {
  value = ["${digitalocean_droplet.test-01.ipv4_address_private }", "${digitalocean_droplet.test-02.ipv4_address_private }"]
}

output "tagged_ips" {
  value = ["${digitalocean_droplet.test-01.ipv4_address_private }"]
}
