provider "google" {
  region      = "${var.region}"
  project     = "${var.project_name}"
  credentials = "${file("${var.credentials_file_path}")}"
}

resource "google_compute_instance" "main" {
  count = 2

  name         = "tf-discover-${count.index}"
  machine_type = "f1-micro"
  zone         = "${var.region_zone}"
  tags         = ["ssh-node", "consul-${count.index}"]

  disk {
    image = "ubuntu-os-cloud/ubuntu-1404-trusty-v20160602"
  }

  network_interface {
    network = "default"

    access_config {
      # Ephemeral
    }
  }

  service_account {
    scopes = ["https://www.googleapis.com/auth/compute.readonly"]
  }
}
