provider "hcloud" {
}

resource "hcloud_server" "test-01" {
  count       = 2
  image       = var.hcloud_image
  name        = "${var.prefix}-01-${count.index + 1}"
  server_type = var.hcloud_size
  location    = var.hcloud_location_nbg
  labels = {
    "${var.prefix}-test-tag" : ""
  }
}

resource "hcloud_server" "test-02" {
  image       = var.hcloud_image
  name        = "${var.prefix}-02"
  server_type = var.hcloud_size
  location    = var.hcloud_location_fsn
}

resource "hcloud_network" "internal" {
  name     = "${var.prefix}-internal"
  ip_range = "10.0.0.0/8"
  labels = {
    "${var.prefix}-test-tag" : ""
  }
}

resource "hcloud_network_subnet" "subnet" {
  network_id   = hcloud_network.internal.id
  type         = "cloud"
  network_zone = "eu-central"
  ip_range     = "10.0.1.0/24"
}

resource "hcloud_server_network" "test-01" {
  count     = 2
  server_id = hcloud_server.test-01[count.index].id
  subnet_id = hcloud_network_subnet.subnet.id
}

resource "hcloud_server_network" "test-02" {
  server_id = hcloud_server.test-02.id
  subnet_id = hcloud_network_subnet.subnet.id
}
