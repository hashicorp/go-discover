variable "name" {}

variable "key_pair_name" {}

variable "private_key_path" {}

variable "subnet_id" {}

variable "vpc_id" {}

variable "tags" {
  type    = "map"
  default = {}
}

variable "username" {
  default = "ubuntu"
}
