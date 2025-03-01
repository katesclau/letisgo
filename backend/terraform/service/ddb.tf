resource "aws_dynamodb_table" "single_table" {
  name         = "${local.prefix}-table"
  billing_mode = "PAY_PER_REQUEST"
  hash_key     = "pk"
  range_key    = "sk"

  attribute {
    name = "pk"
    type = "S"
  }

  attribute {
    name = "sk"
    type = "S"
  }

  tags = {
    Name         = "${local.prefix}-table"
    Environment  = var.Environment
    ServiceName  = var.ServiceName
    Organization = var.Organization
  }

  global_secondary_index {
    name            = "model-index"
    hash_key        = "model"
    projection_type = "ALL"
  }

  attribute {
    name = "model"
    type = "S"
  }
}
