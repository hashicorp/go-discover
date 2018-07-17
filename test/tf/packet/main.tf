provider "packet" {}

resource "packet_project" "project" {
  name = "go-discover:packet-project-${element(random_string.vm_name_suffix.*.result, count.index)}"
}

resource "random_string" "vm_name_suffix" {
  count   = 2
  length  = 8
  upper   = false
  special = false
}

resource "packet_device" "discover-packet01" {
  count            = 2
  hostname         = "go-discover.packet-device-${element(random_string.vm_name_suffix.*.result, count.index)}"
  plan             = "baremetal_1"
  facility         = "ewr1"
  operating_system = "coreos_stable"
  billing_cycle    = "hourly"
  project_id       =  "${packet_project.project.id}"
}

output "project_id" {
  value = "${packet_project.project.id}"
}
