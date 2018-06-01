provider "alicloud" {}

resource "alicloud_instance" "test" {
  image_id        = "centos_7_04_64_20G_alibase_201701015.vhd"
  availability_zone = "us-west-1"
  instance_type   = "ecs.g5.large"
  security_groups = ["${alicloud_security_group.sg.id}"]
  tags            = "${var.tag}"
}

resource "alicloud_security_group" "sg" {
  id          = "default"
  description = "Default"
  vpc_id      = "test_0315"
}
