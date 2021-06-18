
terraform {
  required_version = ">= 0.13"
  required_providers {
    metal = {
      source  = "equinix/metal"
      version = "~> 2.1"
    }
  }
}
