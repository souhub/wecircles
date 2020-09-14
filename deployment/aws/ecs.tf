# _____________________________Task Definitions______________________________________
resource "aws_ecs_task_definition" "task" {
  family                   = "wecircles"
  cpu                      = 256
  memory                   = 512
  network_mode             = "awsvpc"
  requires_compatibilities = ["FARGATE"]
  execution_role_arn       = "arn:aws:iam::535411933495:role/ecsTaskExecutionRole"
  container_definitions    = file("task-definitions.json")
}

# _____________________________Cluster______________________________________
resource "aws_ecs_cluster" "cluster" {
  name = "wecircles"
}

# _____________________________Service______________________________________
resource "aws_ecs_service" "service" {
  name                              = "wecircles"
  cluster                           = aws_ecs_cluster.cluster.arn
  task_definition                   = aws_ecs_task_definition.task.arn
  desired_count                     = 2
  launch_type                       = "FARGATE"
  platform_version                  = "1.4.0"
  health_check_grace_period_seconds = 60

  network_configuration {
    assign_public_ip = true
    security_groups  = [aws_security_group.ec2_sg.id]

    subnets = [
      aws_subnet.publicA.id,
      aws_subnet.publicC.id,
    ]
  }

  load_balancer {
    target_group_arn = aws_lb_target_group.tg.arn
    container_name   = "wecircles" //container_definition.jsonのname
    container_port   = 80
  }
}
