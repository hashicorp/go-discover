provider "digitalocean" {
    token = "${var.digitalocean_token}"
}

resource "digitalocean_tag" "test" {
  name = "${var.prefix}-test-tag"
}

resource "digitalocean_ssh_key" "test" {
  name       = "${var.prefix}-test-key"
  public_key = "${file("${var.ssh_public_path}")}"
}

resource "digitalocean_droplet" "test-01" {
    image              = "${var.do_image}"
    name               = "${var.prefix}-01"
    region             = "${var.do_region}"
    size               = "${var.do_size}"
    private_networking = true
    ssh_keys           = ["${digitalocean_ssh_key.test.id}"]
    tags               = ["${digitalocean_tag.test.id}"]

  provisioner "file" {
    source      = "discover"
    destination = "/tmp/discover"

    connection {
      type        = "ssh"
      user        = "root"
      private_key = "${file("${var.ssh_private_path}")}"
      agent       = false
    }
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "root"
      private_key = "${file("${var.ssh_private_path}")}"
      agent       = false
    }

    inline = [
      "chmod +x /tmp/discover"
    ]
  }
}

resource "digitalocean_droplet" "test-02" {
    image              = "${var.do_image}"
    name               = "${var.prefix}-02"
    region             = "${var.do_region}"
    size               = "${var.do_size}"
    private_networking = true
    ssh_keys           = ["${digitalocean_ssh_key.test.id}"]

  provisioner "file" {
    source      = "discover"
    destination = "/tmp/discover"

    connection {
      type        = "ssh"
      user        = "root"
      private_key = "${file("${var.ssh_private_path}")}"
      agent       = false
    }
  }

  provisioner "remote-exec" {
    connection {
      type        = "ssh"
      user        = "root"
      private_key = "${file("${var.ssh_private_path}")}"
      agent       = false
    }

    inline = [
      "chmod +x /tmp/discover"
    ]
  }
}

resource "digitalocean_firewall" "test" {
  name = "${var.prefix}-test-firewall"

  droplet_ids = ["${digitalocean_droplet.test-01.id }","${digitalocean_droplet.test-02.id }"]

  # Allow inbound SSH access.
  inbound_rule = {
      protocol           = "tcp"
      port_range         = "22"
      source_addresses   = ["0.0.0.0/0"]
  }

  outbound_rule = [
    # Allow outbound connection to https://api.digitalocean.com
    {
        protocol                = "tcp"
        port_range              = "443"
        destination_addresses   = ["0.0.0.0/0"]
    },
    # Allow outbound udp lookups.
    {
        protocol                = "udp"
        port_range              = "53"
        destination_addresses   = ["0.0.0.0/0"]
    }
  ]
}