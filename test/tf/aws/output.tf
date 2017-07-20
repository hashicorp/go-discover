
output "public_ips" {
  value = ["${module.vm01.public_ip}", "${module.vm02.public_ip}"]
}

output "private_ips" {
  value = ["${module.vm01.private_ip}", "${module.vm02.private_ip}"]
}

output "tagged_ips" {
  value = ["${module.vm01.private_ip}"]
}

output "dns_servers" {
  value = [
    "${aws_route53_zone.main.name_servers.0}",
    "${aws_route53_zone.main.name_servers.1}",
    "${aws_route53_zone.main.name_servers.2}",
    "${aws_route53_zone.main.name_servers.3}",
  ]
}

