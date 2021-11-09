resource "random_id" "suffix" {
  byte_length = 4
}

data "google_container_engine_versions" "main" {
  location = var.zone
}

resource "google_container_cluster" "cluster" {
  name               = "consul-k8s-${random_id.suffix.dec}"
  project            = var.project
  initial_node_count = 5
  location           = var.zone
  min_master_version = data.google_container_engine_versions.main.latest_master_version
  node_version       = data.google_container_engine_versions.main.latest_master_version
}

resource "null_resource" "kubeconfig" {
  triggers = {
    cluster = google_container_cluster.cluster.id
  }

  # On creation, we want to setup the kubectl credentials. The easiest way
  # to do this is to shell out to gcloud.
  provisioner "local-exec" {
    command = "gcloud container clusters get-credentials --zone=${var.zone} ${google_container_cluster.cluster.name}"
  }

  # On destroy we want to try to clean up the kubectl credentials. This
  # might fail if the credentials are already cleaned up or something so we
  # want this to continue on failure. Generally, this works just fine since
  # it only operates on local data.
  provisioner "local-exec" {
    when       = destroy
    on_failure = continue
    command    = "rm $HOME/.kube/config"
  }
}

resource "kubernetes_pod" "valid" {
  depends_on = [null_resource.kubeconfig]
  count      = 3

  metadata {
    name = "valid-${count.index}"

    labels = {
      app = "valid"
    }
  }

  spec {
    container {
      image = "nginx:1.7.9"
      name  = "echo"
    }
  }
}

resource "kubernetes_pod" "invalid" {
  depends_on = [null_resource.kubeconfig]
  count      = 2

  metadata {
    name = "invalid-${count.index}"

    labels = {
      app = "invalid"
    }
  }

  spec {
    container {
      image = "nginx:1.7.9"
      name  = "echo"
    }
  }
}
