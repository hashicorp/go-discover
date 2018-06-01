variable "tag" {
  type = "map"

  default = {
    "consul" = "server.test"
  }
}

variable "instance_type" {
  default = "ecs.n4.small"
}
