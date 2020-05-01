provider "linode" {
  version = "~> 1.9.3"
}

provider "random" {
  version = "~> 2.2.1"
}

resource "random_id" "test_id" {
  byte_length = 4
  prefix      = "go-discover-"
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
  region     = element(var.regions, count.index)
  type       = "g6-nanode-1"
  private_ip = true
  # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
  # force an interpolation expression to be interpreted as a list by wrapping it
  # in an extra set of list brackets. That form was supported for compatibility in
  # v0.11, but is no longer supported in Terraform v0.12.
  #
  # If the expression in the following list itself returns a list, remove the
  # brackets to avoid interpretation as a list of lists. If the expression
  # returns a single list item then leave it as-is and remove this TODO comment.
  tags = [element(var.tags, count.index)]
}

