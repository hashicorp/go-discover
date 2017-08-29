# go-discover test for DigitalOcean

This is an integration test setup for go-discover on DigitalOcean.

The test creates two Linux VMs with public and private ip addresses and SSH
access for user `root`. One instance is tagged `go-discover-test-tag`

It then runs the `discover` tool on both VMs to discover the private ip
addresses of all servers with the `go-discover-test-tag` tag. Only the ip
address of the first server should be found for the test to succeeed.

## Run the test

1. Set the environment variable `TF_VAR_digitalocean_token` to contain a
   DigitalOcean API token.
2. Run `make build test destroy` and `tail -f tf.log` in another window
