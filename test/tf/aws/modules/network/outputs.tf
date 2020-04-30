output "subnet_id" {
  value = aws_subnet.internal.id
}

output "vpc_id" {
  value = aws_vpc.main.id
}

