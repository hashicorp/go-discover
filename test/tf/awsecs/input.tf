provider "aws" {
  version = "~> 1.15"
}

variable "prefix" {
  default     = "go-discover"
  description = "prefix of the VPC ?"
}

variable "ecs_cluster" {
  description = "ECS cluster name"
  default     = "go-discover-cluster"
}

variable "private_key_path" {
  default     = "tf_rsa"
  description = "path to the private key used for provisioning"
}

variable "public_key_path" {
  default     = "tf_rsa.pub"
  description = "path to the public key used for provisioning"
}

variable "address_space" {
  default     = "10.0.0.0/16"
  description = "address space of the VPC"
}

variable "subnet_cidr_a" {
  default     = "10.0.1.0/24"
  description = "subnet a of the VPC"
}

variable "subnet_cidr_b" {
  default     = "10.0.2.0/24"
  description = "subnet b of the VPC"
}
