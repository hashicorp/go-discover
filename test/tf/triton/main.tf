provider "triton" {}

data "triton_image" "image" {
  name    = "base-64-lts"
  version = "16.4.1"
}

data "triton_network" "public" {
  name = "Joyent-SDC-Public"
}

resource "triton_machine" "test" {
  package  = "g4-highcpu-128M"
  image    = "${data.triton_image.image.id}"
  networks = ["${data.triton_network.public.id}"]
}

resource "triton_machine" "test_tagged" {
  package  = "g4-highcpu-128M"
  image    = "${data.triton_image.image.id}"
  networks = ["${data.triton_network.public.id}"]

  tags {
    consul-role = "server"
  }
}
