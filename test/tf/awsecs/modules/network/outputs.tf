output "subnet_id_a" {
  value = "${aws_subnet.internal_a.id}"
}

output "subnet_id_b" {
  value = "${aws_subnet.internal_b.id}"
}

output "vpc_id" {
  value = "${aws_vpc.main.id}"
}
