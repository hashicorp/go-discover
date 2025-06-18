module "network" {
  source        = "./modules/network"
  address_space = "${var.address_space}"
  subnet_cidr_a   = "${var.subnet_cidr_a}"
  subnet_cidr_b   = "${var.subnet_cidr_b}"
}
data "aws_region" "current" {}

variable "ecs_amis" {
  type = "map"
  default = {
    "us-east-1" = "ami-5253c32d"
    "us-west-2" = "ami-d2f489aa"
    "us-west-1" = "ami-6b81980b"
    "eu-west-3" = "ami-ca75c4b7"
    "eu-west-2" = "ami-3622cf51"
    "eu-west-1" = "ami-c91624b0"
    "eu-central-1" = "ami-10e6c8fb"
    "ap-northeast-2" = "ami-7c69c112"
    "ap-northeast-1" = "ami-f3f8098c"
    "ap-southeast-2" = "ami-bc04d5de"
    "ap-southeast-1" = "ami-b75a6acb"
    "ca-central-1" = "ami-da6cecbe"
    "ap-south-1" = "ami-c7072aa8"
    "sa-east-1" = "ami-a1e2becd"
    "us-gov-west-1" = "ami-03920462"
  }
}

data "aws_iam_policy_document" "ecs-service-policy" {
    statement {
        actions = ["sts:AssumeRole"]

        principals {
            type        = "Service"
            identifiers = ["ecs.amazonaws.com"]
        }
    }
}

resource "aws_iam_role" "ecs-instance-role" {
    name                = "ecs-instance-role"
    path                = "/"
    assume_role_policy  = "${data.aws_iam_policy_document.ecs-instance-policy.json}"
}

data "aws_iam_policy_document" "ecs-instance-policy" {
    statement {
        actions = ["sts:AssumeRole"]

        principals {
            type        = "Service"
            identifiers = ["ec2.amazonaws.com"]
        }
    }
}

resource "aws_iam_role_policy_attachment" "ecs-instance-role-attachment" {
    role       = "${aws_iam_role.ecs-instance-role.name}"
    policy_arn = "arn:aws:iam::aws:policy/service-role/AmazonEC2ContainerServiceforEC2Role"
}

resource "aws_iam_instance_profile" "ecs-instance-profile" {
    name = "ecs-instance-profile"
    path = "/"
    role = "${aws_iam_role.ecs-instance-role.id}"
    provisioner "local-exec" {
      command = "sleep 10"
    }
}

resource "aws_launch_configuration" "ecs-launch-configuration" {
    name                        = "ecs-launch-configuration"
    image_id                    = "${lookup(var.ecs_amis, data.aws_region.current.name)}"
    instance_type               = "t2.micro"
    iam_instance_profile        = "${aws_iam_instance_profile.ecs-instance-profile.id}"
   
    root_block_device {
      volume_type = "standard"
      volume_size = 8
      delete_on_termination = true
    }

    lifecycle {
      create_before_destroy = true
    }

    associate_public_ip_address = "true"
    user_data                   = <<EOF
#!/bin/bash
echo ECS_CLUSTER=${var.ecs_cluster} >> /etc/ecs/ecs.config
EOF
}

resource "aws_autoscaling_group" "ecs-autoscaling-group" {
    name                        = "ecs-autoscaling-group"
    max_size                    = "2"
    min_size                    = "2"
    desired_capacity            = "2"
    vpc_zone_identifier         = ["${module.network.subnet_id_a}", "${module.network.subnet_id_b}"]
    launch_configuration        = "${aws_launch_configuration.ecs-launch-configuration.name}"
    health_check_type           = "EC2"
}



data "aws_ecs_task_definition" "nginx" {
  task_definition = "${aws_ecs_task_definition.nginx.family}"
}

resource "aws_ecs_task_definition" "nginx" {
    family                = "nginx"
    container_definitions = <<DEFINITION
[
  {
    "name": "nginx",
    "image": "nginx:1.14-alpine",
    "essential": true,
    "portMappings": [
      {
        "hostPort": 0,
        "containerPort": 80
      }
    ],
    "memory": 128
  }
]
DEFINITION
}

resource "aws_ecs_service" "test-ecs-service" {
  	name            = "test-ecs-service"
  	cluster         = "${aws_ecs_cluster.test-ecs-cluster.id}"
  	task_definition = "${aws_ecs_task_definition.nginx.family}:${max("${aws_ecs_task_definition.nginx.revision}")}"
  	desired_count   = 2
}

resource "aws_ecs_cluster" "test-ecs-cluster" {
    name = "${var.ecs_cluster}"
}