resource "aws_key_pair" "main" {
  key_name   = "${var.prefix}-key"
  public_key = "${file(var.public_key_path)}"
}

module "network" {
  source        = "./modules/network"
  address_space = "${var.address_space}"
  subnet_cidr   = "${var.subnet_cidr}"
}

module "vm01" {
  source           = "./modules/virtual_machine"
  name             = "vm01"
  key_pair_name    = "${aws_key_pair.main.key_name}"
  subnet_id        = "${module.network.subnet_id}"
  vpc_id           = "${module.network.vpc_id}"
  private_key_path = "${var.private_key_path}"

  tags {
    "consul" = "server"
  }
}

module "vm02" {
  source           = "./modules/virtual_machine"
  name             = "vm02"
  key_pair_name    = "${aws_key_pair.main.key_name}"
  subnet_id        = "${module.network.subnet_id}"
  vpc_id           = "${module.network.vpc_id}"
  private_key_path = "${var.private_key_path}"
}

