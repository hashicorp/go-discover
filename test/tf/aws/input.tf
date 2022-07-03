provider "aws" {}

variable "prefix" {
  default     = "go-discover"
  description = "prefix of the VPC ?"
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

variable "subnet_cidr" {
  default     = "10.0.1.0/24"
  description = "subnet of the VPC"
}

