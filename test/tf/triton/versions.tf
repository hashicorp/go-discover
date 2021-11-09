
terraform {
  required_version = ">= 0.12"

  required_providers {
    triton = {
      source  = "joyent/triton"
      version = "0.8.2"
    }
  }
}
