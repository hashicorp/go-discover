provider "linode" {
  version = "~> 1.5.0"
}

provider "random" {
  version = "~> 2.0.0"
}

resource "random_id" "test_id" {
  byte_length = 4
  prefix = "go-discover-"
}

variable "regions" {
  default = ["us-east", "ap-south", "ap-south"]
}

variable "tags" {
  default = ["gd-tag1", "gd-tag1", "gd-tag2"]
}

resource "linode_instance" "go-discover-linode" {
  count      = 3
  image      = "linode/debian9"
  label      = "${random_id.test_id.dec}-${count.index}"
  region     = "${element(var.regions, count.index)}"
  type       = "g6-nanode-1"
  private_ip = true
  tags       = ["${element(var.tags, count.index)}"]
}
