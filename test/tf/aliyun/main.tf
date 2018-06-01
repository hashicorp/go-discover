provider "alicloud" {}

resource "alicloud_instance" "test" {
  image_id        = "centos_7_04_64_20G_alibase_201701015.vhd"
  instance_type   = "${var.instance_type}"
  security_groups = ["${alicloud_security_group.default.id}"]
  tags            = "${var.tag}"
}

resource "alicloud_security_group" "default" {
  name        = "default"
  description = "default"
}
