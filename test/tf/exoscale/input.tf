variable "zone" {
  default = "de-fra-1"
}

data "exoscale_compute_template" "ubuntu" {
  zone = var.zone
  name = "Linux Ubuntu 20.04 LTS 64-bit"
}

