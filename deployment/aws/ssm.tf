# _____________________________RDS Password______________________________________
resource "aws_ssm_parameter" "db_pass" {
  name  = "/wecircles/rds/password"
  type  = "String"
  value = aws_db_instance.db.password
}

# _____________________________RDS Endpoint______________________________________
resource "aws_ssm_parameter" "db_endpoint" {
  name  = "/wecircles/rds/endpoint"
  type  = "String"
  value = aws_db_instance.db.endpoint
}
