# Copyright IBM Corp. 2017, 2025
# SPDX-License-Identifier: MPL-2.0

output "subnet_id" {
  value = aws_subnet.internal.id
}

output "vpc_id" {
  value = aws_vpc.main.id
}

