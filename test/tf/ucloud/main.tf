provider "ucloud" {
  version = "~> 1.12.1"
}

resource "ucloud_vpc" "vpc" {
  name        = "UCloud"
  cidr_blocks = [
    "10.0.0.0/16"]
}

resource "ucloud_subnet" "subnet" {
  name       = "UCloud"
  cidr_block = "10.0.0.0/16"
  vpc_id     = "${ucloud_vpc.vpc.id}"
}

data "ucloud_images" "centos" {
  availability_zone = "${var.zone}"
  image_type        = "base"
  name_regex        = "^CentOS 7.[1-7] 64"
  most_recent       = true
}

resource "ucloud_instance" "instance" {
  count             = 2
  availability_zone = "${var.zone}"
  tag               = "UCloud"
  image_id          = "${data.ucloud_images.centos.images[0].id}"
  instance_type     = "n-highcpu-1"
  vpc_id            = "${ucloud_vpc.vpc.id}"
  subnet_id         = "${ucloud_subnet.subnet.id}"
}

resource "ucloud_eip" "eip" {
  count         = 2
  internet_type = "${var.internet_type}"
  charge_mode   = "traffic"
  charge_type   = "dynamic"
}

resource "ucloud_eip_association" "association" {
  count       = 2
  eip_id      = "${ucloud_eip.eip.*.id[count.index]}"
  resource_id = "${ucloud_instance.instance.*.id[count.index]}"
}