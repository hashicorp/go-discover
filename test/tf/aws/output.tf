
output "public_ips" {
  value = ["${module.vm01.public_ip}", "${module.vm02.public_ip}"]
}

output "private_ips" {
  value = ["${module.vm01.private_ip}", "${module.vm02.private_ip}"]
}

output "tagged_ips" {
  value = ["${module.vm01.private_ip}"]
}

