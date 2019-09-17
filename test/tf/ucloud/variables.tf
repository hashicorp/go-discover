variable "zone" {}
// "bgp" for mainland China region and "international" for international region
variable internet_type {
  type    = "string"
  default = "bgp"
}