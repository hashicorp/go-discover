data "aws_ami" "ubuntu" {
  most_recent = true

  filter {
    name   = "name"
    values = ["ubuntu/images/hvm-ssd/ubuntu-trusty-14.04-amd64-server-*"]
  }

  filter {
    name   = "virtualization-type"
    values = ["hvm"]
  }

  owners = ["099720109477"] # Canonical
}

resource "aws_instance" "main" {
  ami                         = "${data.aws_ami.ubuntu.id}"
  instance_type               = "t2.micro"
  key_name                    = "${var.key_pair_name}"
  vpc_security_group_ids      = ["${aws_security_group.ssh.id}"]
  subnet_id                   = "${var.subnet_id}"
  associate_public_ip_address = true

  tags = "${var.tags}"

  connection {
    type        = "ssh"
    host        = "${self.public_ip}"
    user        = "${var.username}"
    private_key = "${file(var.private_key_path)}"
    timeout     = "2m"
  }

  provisioner "file" {
    source      = "discover"
    destination = "/home/${var.username}/discover"
  }

  provisioner "remote-exec" {
    inline = [
      "chmod +x /home/${var.username}/discover",
    ]
  }
}
