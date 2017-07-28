variable "region" {
  default = "europe-west1"
}

variable "region_zone" {
  default = "europe-west1-c"
}

variable "project_name" {
  description = "The ID of the Google Cloud project"
  default     = ""
}

variable "credentials_file_path" {
  description = "Path to the JSON file used to describe your account credentials"
  default     = "tf_gce.json"
}

variable "public_key_path" {
  description = "Path to file containing public key"
  default     = "tf_rsa.pub"
}

variable "private_key_path" {
  description = "Path to file containing private key"
  default     = "tf_rsa"
}

