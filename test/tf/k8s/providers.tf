provider "google" {
  project = var.project
  version = "~> 4.0.0"
}

provider "kubernetes" {
  config_path = "~/.kube/config"
}

provider "local" {
  version = "~> 1.4.0"
}

provider "random" {
  version = "~> 2.2.1"
}

