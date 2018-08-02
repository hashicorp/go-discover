provider "google" {
  project = "${var.project}"
  version = "~> 1.16"
}

provider "kubernetes" {
  host = "https://${google_container_cluster.cluster.endpoint}"
  client_certificate = "${base64decode(google_container_cluster.cluster.master_auth.0.client_certificate)}"
  client_key = "${base64decode(google_container_cluster.cluster.master_auth.0.client_key)}"
  cluster_ca_certificate = "${base64decode(google_container_cluster.cluster.master_auth.0.cluster_ca_certificate)}"
  load_config_file = false
}

provider "local" {
  version = "~> 1.0"
}

provider "random" {
  version = "~> 1.0"
}
