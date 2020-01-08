data "oci_identity_availability_domains" "ads" {
  compartment_id = "${var.tenancy_ocid}"
}

variable "tenancy_ocid" {
  description = "The Tenancy OCID for the target account."
}

variable "user_ocid" {
  description = "The User OCID for the target account credentials."
}

variable "fingerprint" {
  description = "The fingerprint for the key being used to connect to OCI."
}

variable "private_key_file" {
  description = "The path to private key to authenticate the connection OCI."
}

variable "region" {
  description = "The region in which to create the resources."
}

variable "image_map" {
  type        = "map"
  description = "A map of instances keyed by the region."
  default     = {
    ap-seoul-1     = "ocid1.image.oc1.ap-seoul-1.aaaaaaaayajtmksg4tot2pvrezgmqbbhgul5co5flnfvx6avt23hvcdtnk3a"
    ap-tokyo-1     = "ocid1.image.oc1.ap-tokyo-1.aaaaaaaa7ggytzvqrgjaxgylzpy4u64puuml2yjfhys4m2thznuwygdyxzzq"
    ca-toronto-1   = "ocid1.image.oc1.ca-toronto-1.aaaaaaaanyidznndvmpfwv2uybfotbgr7rm6v5rhrecltptx2dxg76d6gdva"
    eu-frankfurt-1 = "ocid1.image.oc1.eu-frankfurt-1.aaaaaaaap4e7y2fyzyx57bfxg6rs5zbnbepfmvcjkezfnhnb4tjo77hl2cma"
    uk-london-1    = "ocid1.image.oc1.uk-london-1.aaaaaaaaoy6vxajijvdxwr432alpjqfxokuuserhj7qof6vfv53o2vvxahrq"
    us-ashburn-1   = "ocid1.image.oc1.iad.aaaaaaaa64ahfqwfhk7ft53o2vc4gz2hb7tiugfjxxaafejbin4zjbg3anpq"
    us-phoenix-1   = "ocid1.image.oc1.phx.aaaaaaaa54bjroabxqmvzfrpq7vyp4yoga33dyfbtdaufdr2dkhubnh4szyq"
  }
}

variable "tag_namespace" {
  description = "The tag namespace to use for testing. Defaults to defined."
  default     = "defined"
}

variable "tag_key" {
  description = "The tag key to use for testing. Defaults to discover."
  default     = "discover"
}

variable "tag_value" {
  description = "The tag value to use for testing. Defaults to me."
  default     = "me"
}
