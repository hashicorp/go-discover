# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

output "subnet_id" {
  value = aws_subnet.internal.id
}

output "vpc_id" {
  value = aws_vpc.main.id
}

