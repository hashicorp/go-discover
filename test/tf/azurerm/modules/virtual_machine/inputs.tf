variable "name" {
}

variable "resource_group_name" {
}

variable "location" {
}

variable "subnet_id" {
}

variable "size" {
  default = "Standard_A1_v2"
}

variable "username" {
  default = "ubuntu"
}

variable "tags" {
  type    = map(string)
  default = {}
}

