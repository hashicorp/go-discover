provider "metal" {
  version = "~> 2.8.1"
}

provider "random" {
  version = "~> 2.2.1"
}

variable "metro" {
  default = ["sv"]
}

variable "tags" {
  default = ["tag1", "tag1", "tag2", "tag3"]
}

variable "metal_project" {
  description = "Existing Equinix Metal project"
}

resource "random_string" "vm_name_suffix" {
  count   = 4
  length  = 8
  upper   = false
  special = false
}

resource "metal_device" "discover-metal01" {
  count    = 4
  hostname = "go-discover.metal-device-${element(random_string.vm_name_suffix.*.result, count.index)}"
  plan     = "c3.small.x86"
  metro    = var.metro
  # TF-UPGRADE-TODO: In Terraform v0.10 and earlier, it was sometimes necessary to
  # force an interpolation expression to be interpreted as a list by wrapping it
  # in an extra set of list brackets. That form was supported for compatibility in
  # v0.11, but is no longer supported in Terraform v0.12.
  #
  # If the expression in the following list itself returns a list, remove the
  # brackets to avoid interpretation as a list of lists. If the expression
  # returns a single list item then leave it as-is and remove this TODO comment.
  tags             = [element(var.tags, count.index)]
  operating_system = "ubuntu_18_04"
  billing_cycle    = "hourly"
  project_id       = var.metal_project
}

