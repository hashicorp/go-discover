provider "packet" {}

resource "packet_project" "project" {
  name = "go-discover:packet"
}

resource "packet_device" "discover-packet" {
  count            = 2
  hostname         = "tf.discover"
  plan             = "baremetal_1"
  facility         = "ewr1"
  operating_system = "coreos_stable"
  billing_cycle    = "hourly"
  project_id       = "${packet_project.project.id}"
}

output "project_id" {
  value = "${packet_project.project.id}"
}
