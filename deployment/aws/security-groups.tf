# _____________________________Load Balancer______________________________________
resource "aws_security_group" "lb_sg" {
  name        = "lb"
  description = "From Internet to LB"
  vpc_id      = aws_vpc.vpc.id

  ingress {
    description = "HTTPS from Internet"
    from_port   = 443
    to_port     = 443
    protocol    = "tcp"
    cidr_blocks = ["0.0.0.0/0"]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# _____________________________EC2______________________________________
resource "aws_security_group" "ec2_sg" {
  name        = "ec2"
  description = "From LB to EC2"
  vpc_id      = aws_vpc.vpc.id

  ingress {
    description     = "HTTP from Load Balancer"
    from_port       = 80
    to_port         = 80
    protocol        = "tcp"
    security_groups = [aws_security_group.lb_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}

# _____________________________RDS______________________________________
resource "aws_security_group" "db_sg" {
  name        = "db"
  description = "From EC2 to MySQL"
  vpc_id      = aws_vpc.vpc.id

  ingress {
    description     = "MySQL from EC2"
    from_port       = 3306
    to_port         = 3306
    protocol        = "tcp"
    security_groups = [aws_security_group.ec2_sg.id]
  }

  egress {
    from_port   = 0
    to_port     = 0
    protocol    = "-1"
    cidr_blocks = ["0.0.0.0/0"]
  }
}
