resource "aws_s3_bucket" "public_bucket" {
  bucket = "${local.prefix}-public-bucket"

  tags = {
    Name         = "${local.prefix}-public-bucket"
    Environment  = var.Environment
    Organization = var.Organization
  }
}
