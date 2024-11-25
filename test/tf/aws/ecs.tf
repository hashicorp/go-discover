# Copyright (c) HashiCorp, Inc.
# SPDX-License-Identifier: MPL-2.0

data "aws_iam_policy" "ecs_execution_role_managed_policy" {
  name = "AmazonECSTaskExecutionRolePolicy"
}

data "aws_iam_policy_document" "ecs_execution_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_execution_role" {
  assume_role_policy = data.aws_iam_policy_document.ecs_execution_role_policy.json

  managed_policy_arns = [
    data.aws_iam_policy.ecs_execution_role_managed_policy.arn
  ]
}

data "aws_iam_policy_document" "ecs_task_role_policy" {
  statement {
    actions = ["sts:AssumeRole"]

    principals {
      type        = "Service"
      identifiers = ["ecs-tasks.amazonaws.com"]
    }
  }
}

resource "aws_iam_role" "ecs_task_role" {
  assume_role_policy = data.aws_iam_policy_document.ecs_task_role_policy.json
}

data "aws_iam_policy_document" "ecs_auto_discover" {
  statement {
    effect = "Allow"
    actions = [
      "ecs:ListClusters",
      "ecs:ListServices",
      "ecs:DescribeServices",
      "ecs:ListTasks",
      "ecs:DescribeTasks",
    ]
    resources = ["*"]
  }
}

resource "aws_iam_role_policy" "ecs_auto_discover" {
  policy = data.aws_iam_policy_document.ecs_auto_discover.json
  role   = aws_iam_role.ecs_task_role.id
}

resource "aws_ecs_cluster" "ecs1" {
  name = "${var.prefix}-1"
}

resource "aws_ecs_cluster" "ecs2" {
  name = "${var.prefix}-2"
}

resource "aws_ecs_task_definition" "task_def" {
  family                   = var.prefix
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  cpu                      = 256
  memory                   = 512

  container_definitions = jsonencode([
    {
      name : var.prefix
      image : "public.ecr.aws/docker/library/busybox:stable"
      cpu : 256
      memory : 512
      essential : true
      entrypoint: [ "/bin/sh", "-c" ]
      command: [ "while true; do sleep 30; done;" ]
    }
  ])
}

resource "aws_ecs_task_definition" "task_def_familia" {
  family                   = "${var.prefix}-familia"
  requires_compatibilities = ["FARGATE"]
  network_mode             = "awsvpc"
  execution_role_arn       = aws_iam_role.ecs_execution_role.arn
  task_role_arn            = aws_iam_role.ecs_task_role.arn
  cpu                      = 256
  memory                   = 512

  container_definitions = jsonencode([
    {
      name : var.prefix
      image : "public.ecr.aws/docker/library/busybox:stable"
      cpu : 256
      memory : 512
      essential : true
      entrypoint: [ "/bin/sh", "-c" ]
      command: [ "while true; do sleep 30; done;" ]
    }
  ])
}

resource "aws_ecs_service" "task_service_tagged" {
  name            = "${var.prefix}-tagged"
  cluster         = aws_ecs_cluster.ecs1.name
  task_definition = aws_ecs_task_definition.task_def.arn
  propagate_tags  = "SERVICE"
  launch_type     = "FARGATE"
  desired_count   = 1

  network_configuration {
    assign_public_ip = true
    subnets          = [module.network.subnet_id]
  }

  tags = {
    consul : "server"
  }
}

resource "aws_ecs_service" "task_service_not_tagged" {
  name            = "${var.prefix}-not-tagged"
  cluster         = aws_ecs_cluster.ecs1.name
  task_definition = aws_ecs_task_definition.task_def.arn
  propagate_tags  = "SERVICE"
  launch_type     = "FARGATE"
  desired_count   = 1

  network_configuration {
    assign_public_ip = true
    subnets          = [module.network.subnet_id]
  }
}

resource "aws_ecs_service" "task_service_familia_tagged" {
  name            = "${var.prefix}-familia"
  cluster         = aws_ecs_cluster.ecs2.name
  task_definition = aws_ecs_task_definition.task_def_familia.arn
  propagate_tags  = "SERVICE"
  launch_type     = "FARGATE"
  desired_count   = 2

  network_configuration {
    assign_public_ip = true
    subnets          = [module.network.subnet_id]
  }

  tags = {
    consul : "server"
  }
}

resource "aws_ecs_service" "task_service_familia_not_tagged" {
  name            = "${var.prefix}-familia-not-tagged"
  cluster         = aws_ecs_cluster.ecs2.name
  task_definition = aws_ecs_task_definition.task_def_familia.arn
  propagate_tags  = "SERVICE"
  launch_type     = "FARGATE"
  desired_count   = 1

  network_configuration {
    assign_public_ip = true
    subnets          = [module.network.subnet_id]
  }
}
