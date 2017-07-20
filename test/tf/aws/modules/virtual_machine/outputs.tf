output "public_ip" {
  value = "${aws_instance.main.public_ip}"
}

output "private_ip" {
  # value = "${aws_network_interface.internal.private_ips[0]}"
  value = "${aws_instance.main.private_ip}"
}
