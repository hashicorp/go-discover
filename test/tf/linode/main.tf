provider "linode" {
  version = "~> 1.5.0"
}

resource "linode_instance" "tagged" {
  count              = 2
  image              = "${var.linode_image}"
  label              = "${var.prefix}-tagged-${count.index}"
  region             = "${var.linode_region}"
  type               = "${var.linode_type}"
  private_ip         = true
  tags               = ["${var.linode_tag}"]
}

resource "linode_instance" "untagged" {
  image              = "${var.linode_image}"
  label              = "${var.prefix}-untagged"
  region             = "${var.linode_region}"
  type               = "${var.linode_type}"
  private_ip         = true
}
