provider "packet" {}

variable "facility" {
  default = ["ewr1", "sjc1", "ams1", "nrt1"]
}

variable "tags" {
  default = ["tag1", "tag1", "tag2", "tag3"]
}

resource "packet_project" "project" {
  name = "go-discover:packet-project-${element(random_string.vm_name_suffix.*.result, count.index)}"
}

resource "random_string" "vm_name_suffix" {
  count   = 4
  length  = 8
  upper   = false
  special = false
}

resource "packet_device" "discover-packet01" {
  count            = 4
  hostname         = "go-discover.packet-device-${element(random_string.vm_name_suffix.*.result, count.index)}"
  plan             = "baremetal_0"
  facility         = "${element(var.facility, count.index)}"
  tags             = ["${element(var.tags, count.index)}"]
  operating_system = "ubuntu_16_04"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.project.id}"
}

output "project_id" {
  value = "${packet_project.project.id}"
}
