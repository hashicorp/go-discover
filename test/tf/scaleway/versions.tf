
terraform {
  required_version = ">= 0.12"

  required_providers {
    scaleway = {
      source  = "scaleway/scaleway"
      version = "2.1.0"
    }
  }
}
