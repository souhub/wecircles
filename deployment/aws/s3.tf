# _____________________________LB Bucket______________________________________
resource "aws_s3_bucket" "lb_logs" {
  bucket        = "wecircles-logs" # Have yo be obly one in the world
  force_destroy = true             # Force destroy

  lifecycle_rule {
    enabled = true

    expiration {
      days = "30"
    }
  }
}

# _____________________________LB Bucket Policy______________________________________
resource "aws_s3_bucket_policy" "lb_logs" {
  bucket = aws_s3_bucket.lb_logs.id
  policy = data.aws_iam_policy_document.lb_logs.json
}

data "aws_iam_policy_document" "lb_logs" {
  statement {
    effect    = "Allow"
    actions   = ["s3:PutObject"]
    resources = ["arn:aws:s3:::${aws_s3_bucket.lb_logs.id}/*"]

    principals {
      type        = "AWS"
      identifiers = ["582318560864"]
    }
  }
}


