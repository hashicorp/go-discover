provider "oci" {
  tenancy_ocid      = "${var.tenancy_ocid}"
  user_ocid         = "${var.user_ocid}"
  fingerprint       = "${var.fingerprint}"
  private_key_path  = "${var.private_key_file}"
  region            = "${var.region}"
}

# Virtual Cloud Network
resource "oci_core_vcn" "vcn" {
  cidr_block     = "10.0.0.0/16"
  compartment_id = "${var.tenancy_ocid}"
}

# Private subnet
resource "oci_core_subnet" "private" {
  cidr_block                 = "10.0.0.0/24"
  compartment_id             = "${var.tenancy_ocid}"
  vcn_id                     = "${oci_core_vcn.vcn.id}"
  prohibit_public_ip_on_vnic = true
}

# Public subnet
resource "oci_core_subnet" "public" {
  cidr_block                 = "10.0.1.0/24"
  compartment_id             = "${var.tenancy_ocid}"
  vcn_id                     = "${oci_core_vcn.vcn.id}"
}

# Private instance with a freeform tag
resource "oci_core_instance" "private_freeform" {
  availability_domain = "${lookup(data.oci_identity_availability_domains.ads.availability_domains[0], "name")}"
  compartment_id      = "${var.tenancy_ocid}"
  shape               = "VM.Standard2.1"
  freeform_tags       = "${map("${var.tag_key}", "${var.tag_value}")}"

  create_vnic_details {
    subnet_id        = "${oci_core_subnet.private.id}"
    assign_public_ip = false
  }

  source_details {
    source_id   = "${var.image_map[var.region]}"
    source_type = "image"
  }
}

# Tag Namespace
resource "oci_identity_tag_namespace" "defined" {
  compartment_id = "${var.tenancy_ocid}"
  description    = "For automated testing of go-discover"
  name           = "${var.tag_namespace}"
}

resource "oci_identity_tag" "tag_key" {
  description      = "For automated testing of go-discover"
  name             = "${var.tag_key}"
  tag_namespace_id = "${oci_identity_tag_namespace.defined.id}"
}


# Private instance with a defined tag
resource "oci_core_instance" "private_defined" {
  availability_domain = "${lookup(data.oci_identity_availability_domains.ads.availability_domains[1], "name")}"
  compartment_id      = "${var.tenancy_ocid}"
  shape               = "VM.Standard2.1"
  defined_tags        = "${map("${var.tag_namespace}.${var.tag_key}", "${var.tag_value}")}"

  create_vnic_details {
    subnet_id        = "${oci_core_subnet.private.id}"
    assign_public_ip = false
  }

  source_details {
    source_id   = "${var.image_map[var.region]}"
    source_type = "image"
  }
}

# Public instance
resource "oci_core_instance" "public" {
  availability_domain = "${lookup(data.oci_identity_availability_domains.ads.availability_domains[2], "name")}"
  compartment_id      = "${var.tenancy_ocid}"
  shape               = "VM.Standard2.1"
  defined_tags        = "${map("${var.tag_namespace}.${var.tag_key}", "${var.tag_value}")}"

  create_vnic_details {
    subnet_id        = "${oci_core_subnet.public.id}"
    assign_public_ip = true
  }

  source_details {
    source_id   = "${var.image_map[var.region]}"
    source_type = "image"
  }
}