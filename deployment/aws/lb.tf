# _____________________________Load Balancer______________________________________
resource "aws_lb" "lb" {
  name               = "wecircles"
  load_balancer_type = "application"
  security_groups    = [aws_security_group.lb_sg.id]
  subnets            = [aws_subnet.publicA.id, aws_subnet.publicC.id]

  access_logs {
    bucket  = aws_s3_bucket.lb_logs.bucket
    prefix  = "wecircles"
    enabled = true
  }
}

# _____________________________Lisner______________________________________
resource "aws_lb_listener" "lisner" {
  load_balancer_arn = aws_lb.lb.arn
  port              = "443"
  protocol          = "HTTPS"
  ssl_policy        = "ELBSecurityPolicy-2016-08"
  certificate_arn   = "arn:aws:acm:ap-northeast-1:535411933495:certificate/3b7a873b-7942-4b1d-9dc9-6b74e01fd1cb"

  default_action {
    type             = "forward"
    target_group_arn = aws_lb_target_group.tg.arn
  }
}

# _____________________________Target Group______________________________________
resource "aws_lb_target_group" "tg" {
  name        = "wecircles"
  vpc_id      = aws_vpc.vpc.id
  target_type = "ip"
  port        = 80
  protocol    = "HTTP"

  health_check {
    path                = "/"
    healthy_threshold   = 5
    unhealthy_threshold = 2
    timeout             = 5
    interval            = 30
    matcher             = 200
    port                = "traffic-port"
    protocol            = "HTTP"
  }

  depends_on = [aws_lb.lb]

}
