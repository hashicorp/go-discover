//Configuring the provider
provider "exoscale" {
  version = "~> 0.15"
}

resource "exoscale_compute" "test-01" {
  display_name = "test-01"
  disk_size = 10
  size = "Tiny"
  template_id = data.exoscale_compute_template.ubuntu.id
  zone = var.zone

  tags = {
    test-01 = "consul"
  }
}

resource "exoscale_compute" "test-02" {
  display_name = "test-02"
  disk_size = 10
  size = "Tiny"
  template_id = data.exoscale_compute_template.ubuntu.id
  zone = "ch-gva-2"

  ip6 = true

  tags = {
    test-02 = "consul"
  }
}
