provider "vultr" {
  api_key = ""
}

resource "vultr_server" "test" {
  count			= 2
  plan_id		= "201"
  region_id		= "1"
  os_id			= "167"
  label			= "go-discover-test-${count.index}"
  tag			= "go-discover-test-tag"
  enable_private_network = true
}
