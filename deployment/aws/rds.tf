# _____________________________Parameter Group______________________________________
resource "aws_db_parameter_group" "parameter_group" {
  name   = "wecircles-parameters"
  family = "mysql8.0"

  parameter {
    name  = "character_set_database"
    value = "utf8mb4"
  }

  parameter {
    name  = "character_set_server"
    value = "utf8mb4"
  }

}

# _____________________________Option Group______________________________________
resource "aws_db_option_group" "option_group" {
  name                 = "wecircles-options"
  engine_name          = "mysql"
  major_engine_version = "8.0"
}

# _____________________________Subnet Group______________________________________
resource "aws_db_subnet_group" "subnet_group" {
  name       = "wecircles-subnets"
  subnet_ids = [aws_subnet.privateA.id, aws_subnet.privateC.id]
}

# _____________________________Instance______________________________________
resource "aws_db_instance" "db" {
  identifier                 = "wecircles"
  name                       = "wecircles"
  engine                     = "mysql"
  engine_version             = "8.0.20"
  instance_class             = "db.t2.micro"
  allocated_storage          = 20
  max_allocated_storage      = 1000
  storage_type               = "gp2"
  username                   = "root"
  password                   = random_password.password.result
  vpc_security_group_ids     = [aws_security_group.db_sg.id]
  multi_az                   = false
  publicly_accessible        = false
  auto_minor_version_upgrade = false
  skip_final_snapshot        = true
  port                       = 3306
  apply_immediately          = true
  parameter_group_name       = aws_db_parameter_group.parameter_group.name
  option_group_name          = aws_db_option_group.option_group.name
  db_subnet_group_name       = aws_db_subnet_group.subnet_group.name
}

# _____________________________Password Generation______________________________________
resource "random_password" "password" {
  length  = 16
  special = false
}




