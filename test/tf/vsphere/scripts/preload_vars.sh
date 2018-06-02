#!/usr/bin/env bash

# preload_vars.sh prompts you for environment variables that need to be set to
# ensure the build runs properly. Note that we don't do any sanity checking -
# if your build fails, re-run this script to attempt to add the correct vars.
#
# In addition to helping facilitate manual execution of the acceptance tests,
# this also helps to document the variables necessary for testing.
#
# Make sure you source this script, versus execute it:
#   source ./preload_vars.sh

# Connectivity
vsphere_host=""
vsphere_username=""
vsphere_password=""
skip_ssl=""

# Terraform vars
datacenter=""
datastore=""
cluster=""
network_label=""
template_name=""

read -r -p "Enter your vCenter/ESXi host: " vsphere_host
read -r -p "Enter your vCenter/ESXi username: " vsphere_username
read -r -s -p "Enter your vCenter/ESXi password: " vsphere_password
# newline is needed here because it is not echoed on slient input
echo
read -r -p "Do you want to skip SSL validation? (Please type \"yes\" if desired) " skip_ssl

read -r -p "Enter the vSphere datacenter to connect to: " datacenter
read -r -p "Enter the vSphere datastore to use: " datastore
read -r -p "Enter the vSphere cluster to use: " cluster
read -r -p "Enter the network name to use for VMs: " network_label
read -r -p "Enter the VM template name: " template_name

export VSPHERE_SERVER="${vsphere_host}"
export VSPHERE_USER="${vsphere_username}"
export VSPHERE_PASSWORD="${vsphere_password}"
if [ "${skip_ssl}" == "yes" ]; then
  export VSPHERE_ALLOW_UNVERIFIED_SSL="true"
fi

export TF_VAR_datacenter="${datacenter}"
export TF_VAR_datastore="${datastore}"
export TF_VAR_cluster="${cluster}"
export TF_VAR_network_label="${network_label}"
export TF_VAR_template_name="${template_name}"
