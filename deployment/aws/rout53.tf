# Association between LB and Domain
resource "aws_route53_record" "wecircles" {
  zone_id = "Z03094043IH28G2M7NJEE"
  name    = "wecircles.net"
  type    = "A"

  alias {
    name                   = aws_lb.lb.dns_name
    zone_id                = aws_lb.lb.zone_id
    evaluate_target_health = true
  }

}
