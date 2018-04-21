provider "digitalocean" {
  token = "${var.digitalocean_token}"
}

resource "digitalocean_tag" "test" {
  name = "${var.prefix}-test-tag"
}

resource "digitalocean_droplet" "test-01" {
  image              = "${var.do_image}"
  name               = "${var.prefix}-01"
  region             = "${var.do_region}"
  size               = "${var.do_size}"
  private_networking = true
  tags               = ["${digitalocean_tag.test.id}"]
}

resource "digitalocean_droplet" "test-02" {
  image              = "${var.do_image}"
  name               = "${var.prefix}-02"
  region             = "${var.do_region}"
  size               = "${var.do_size}"
  private_networking = true
}

resource "digitalocean_firewall" "test" {
  name = "${var.prefix}-test-firewall"

  droplet_ids = ["${digitalocean_droplet.test-01.id }", "${digitalocean_droplet.test-02.id }"]
}
