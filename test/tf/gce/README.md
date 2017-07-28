# go-discover test for Google Cloud

This is an integration test setup for go-discover on Google Cloud.

The test creates two Linux VMs with public and private ip addresses and SSH
access for user `ubuntu`. The instances are tagged with `consul-0` and
`consul-1` respectively.

It then runs the `discover` tool on both VMs to discover the private ip
addresses of all servers with the `consul-0` tag. Only the ip address of the
first server should be found for the test to succeeed.

## Create a GCE Project and Credentials

1. Go to https://console.cloud.google.com/
1. IAM &amp; Admin / Settings: 
	* Create Project, e.g. `discover`
	* Write down the `Project ID`, e.g. `discover-xxx`
1. Billing: Ensure that the project is linked to a billing account
1. API Manager / Dashboard: Enable the following APIs
	* Google Compute Engine API
1. IAM &amp; Admin / Service Accounts: Create Service Account
	* Service account name: `admin`
	* Roles:
		* `Project/Service Account Actor`
		* `Compute Engine/Compute Instance Admin (v1)`
		* `Compute Engine/Compute Security Admin`
	* Furnish a new private key: `yes`
	* Key type: `JSON`
1. The credentials file `discover-xxx.json` will have been downloaded
   automatically to your machine

## Run the test

1. Update the `project_id` and the `credentials_file_path` in the `vars.tf` file
1. Run `make build test destroy` and `tail -f tf.log` in another window

