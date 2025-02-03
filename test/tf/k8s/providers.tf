# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

provider "google" {
  project = var.project
  version = "~> 3.19.0"
}

provider "kubernetes" {
  version = "~> 1.11.1"
  host    = "https://${google_container_cluster.cluster.endpoint}"
  cluster_ca_certificate = base64decode(
    google_container_cluster.cluster.master_auth[0].cluster_ca_certificate,
  )
  username         = google_container_cluster.cluster.master_auth[0].username
  password         = google_container_cluster.cluster.master_auth[0].password
  load_config_file = false
}

provider "local" {
  version = "~> 1.4.0"
}

provider "random" {
  version = "~> 2.2.1"
}

