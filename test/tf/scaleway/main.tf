//Configuring the provider
provider "scaleway" {
  version = "~> 1.8.0"
  region  = "${var.region}"
}

resource "scaleway_server" "test" {
  count = 2
  image = "${var.image}"
  type  = "C2S"
  tags  = ["consul-server"]
  name  = "test-server"
}

resource "scaleway_server" "dummy" {
  image = "${var.image}"
  name  = "dummy_scaleway_instance"
  type  = "C2S"
  tags  = ["Dummy"]
}
