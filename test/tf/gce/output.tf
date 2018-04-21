output "public_ips" {
  value = "${join(" ", google_compute_instance.main.*.network_interface.0.access_config.0.assigned_nat_ip)}"
}

output "private_ips" {
  value = "${join(" ", google_compute_instance.main.*.network_interface.0.address)}"
}

output "tagged_ips" {
  value = "${google_compute_instance.main.0.network_interface.0.address}"
}
