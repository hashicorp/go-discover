provider "packet" {
  version = "~> 2.8.1"
}

provider "random" {
  version = "~> 2.2.1"
}

variable "facility" {
  default = ["ewr1", "sjc1", "ams1", "nrt1"]
}

variable "tags" {
  default = ["tag1", "tag1", "tag2", "tag3"]
}

variable "packet_project" {
  description = "Existing packet project"
}

resource "random_string" "vm_name_suffix" {
  count   = 4
  length  = 8
  upper   = false
  special = false
}

resource "packet_device" "discover-packet01" {
  count      = 4
  hostname   = "go-discover.packet-device-${element(random_string.vm_name_suffix.*.result, count.index)}"
  plan       = "baremetal_0"
  facilities = [var.facility[count.index]]
  # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
  # force an interpolation expression to be interpreted as a list by wrapping it
  # in an extra set of list brackets. That form was supported for compatibility in
  # v0.11, but is no longer supported in Terraform v0.12.
  #
  # If the expression in the following list itself returns a list, remove the
  # brackets to avoid interpretation as a list of lists. If the expression
  # returns a single list item then leave it as-is and remove this TODO comment.
  tags             = [element(var.tags, count.index)]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = var.packet_project
}

