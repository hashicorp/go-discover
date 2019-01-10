provider "alicloud" {
  version = "1.24.0"
}

resource "alicloud_instance" "test" {
  count           = 2
  image_id        = "${var.image_id}"
  instance_type   = "${var.instance_type}"
  security_groups = ["${alicloud_security_group.default.id}"]
  tags            = "${var.tag}"
}

resource "alicloud_security_group" "default" {
  name        = "default"
  description = "default"
}
