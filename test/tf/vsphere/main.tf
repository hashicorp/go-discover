// This is the Terraform file that is necessary to stand up the test
// environment for vSphere node discovery.
//
// In order for this to function correctly, you need to have set up the correct
// vSphere environment variables for connecting. The acceptance tests in
// provider/vsphere also expect these.
//
// The expectant environment variables are:
//
// VSPHERE_SERVER:               The vCenter host
// VSPHERE_USER:                 The vSphere username
// VSPHERE_PASSWORD:             The password for VSPHERE_USER
// VSPHERE_ALLOW_UNVERIFIED_SSL: Set if you need to skip SSL validation
//
// In addition to the above, you need to supply all of the variables below. If
// you are automating in CI, the best way to supply these would be via a
// TF_VAR_ variable, such as TF_VAR_template_name.

// The name of the template to clone from. This needs to be a template
// customizable by VMware tools, such as an up-to-date Linux VM.
variable "template_name" {
  type = string
}

// The name of the datacenter to deploy to.
variable "datacenter" {
  type = string
}

// The name of the datastore to place the VMs in.
variable "datastore" {
  type = string
}

// The name of the cluster to deploy to.
variable "cluster" {
  type = string
}

// The name of the network to deploy to.
variable "network_label" {
  type = string
}

// The number of virtual machines to deploy. For a reasonable acceptance test
// this value should be at least 2.
variable "vm_count" {
  default = "2"
}

data "vsphere_datacenter" "dc" {
  name = var.datacenter
}

data "vsphere_datastore" "datastore" {
  name          = var.datastore
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_compute_cluster" "cluster" {
  name          = var.cluster
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_network" "network" {
  name          = var.network_label
  datacenter_id = data.vsphere_datacenter.dc.id
}

data "vsphere_virtual_machine" "template" {
  name          = var.template_name
  datacenter_id = data.vsphere_datacenter.dc.id
}

resource "vsphere_tag_category" "category" {
  name        = "go-discover-test-category"
  cardinality = "SINGLE"

  associable_types = [
    "VirtualMachine",
  ]
}

resource "vsphere_tag" "tag" {
  name        = "go-discover-test-tag"
  category_id = vsphere_tag_category.category.id
}

resource "random_string" "vm_name_suffix" {
  count   = var.vm_count
  length  = 4
  upper   = false
  special = false
}

resource "vsphere_virtual_machine" "vm" {
  count            = var.vm_count
  name             = "go-discover-test-${element(random_string.vm_name_suffix.*.result, count.index)}"
  resource_pool_id = data.vsphere_compute_cluster.cluster.resource_pool_id
  datastore_id     = data.vsphere_datastore.datastore.id

  num_cpus = 2
  memory   = 1024
  guest_id = data.vsphere_virtual_machine.template.guest_id

  scsi_type = data.vsphere_virtual_machine.template.scsi_type

  network_interface {
    network_id   = data.vsphere_network.network.id
    adapter_type = data.vsphere_virtual_machine.template.network_interface_types[0]
  }

  disk {
    label            = "disk0"
    size             = data.vsphere_virtual_machine.template.disks[0].size
    eagerly_scrub    = data.vsphere_virtual_machine.template.disks[0].eagerly_scrub
    thin_provisioned = data.vsphere_virtual_machine.template.disks[0].thin_provisioned
  }

  clone {
    template_uuid = data.vsphere_virtual_machine.template.id

    customize {
      linux_options {
        host_name = "go-discover-test-${element(random_string.vm_name_suffix.*.result, count.index)}"
        domain    = "test.internal"
      }

      network_interface {
        ipv4_address = cidrhost("10.0.0.0/24", 10 + count.index)
        ipv4_netmask = 24
      }

      ipv4_gateway = "10.0.0.1"
    }
  }

  tags = [vsphere_tag.tag.id]
}

