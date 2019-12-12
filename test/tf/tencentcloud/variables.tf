variable "availability_zone" {
  default = "ap-guangzhou-3"
}

variable "image_id" {
  default = "img-9qabwvbn"
}

variable "instance_type" {
  default = "S1.SMALL1"
}

variable "tag" {
  type = "map"

  default = {
    "consul" = "test"
  }
}
