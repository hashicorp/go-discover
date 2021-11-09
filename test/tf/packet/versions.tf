
terraform {
  required_version = ">= 0.12"

  required_providers {
    packet = {
      source  = "packethost/packet"
      version = "3.2.1"
    }
  }
}
