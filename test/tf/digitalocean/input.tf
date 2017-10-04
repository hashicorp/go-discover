variable "digitalocean_token" {}

variable "prefix" {
  default = "go-discover"
}

variable "do_image" {
  default = "ubuntu-16-04-x64"
}

variable "do_region" {
  default = "nyc3"
}

variable "do_size" {
  default = "512mb"
}

variable "ssh_public_path" {
  default = "./tf_rsa.pub"
}

variable "ssh_private_path" {
  default = "./tf_rsa"
}