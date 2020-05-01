variable "project" {
  default = "go-discover-tests"

  description = <<EOF
Google Cloud Project to launch resources in. This project must have GKE
enabled and billing activated.
EOF

}

variable "zone" {
  default     = "us-east1-b"
  description = "The zone to launch all the GKE nodes in."
}

