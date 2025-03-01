# Defines the Redis cluster, log group, Kinesis Firehose delivery stream, IAM role, and S3 bucket for logs.
# This is used to store sessions, cache data, and is used as the eventbus for our CQRS events.

resource "aws_elasticache_cluster" "redis_cluster" {
  cluster_id           = "${local.prefix}-redis-cluster"
  engine               = "redis"
  node_type            = "cache.t1.micro"
  num_cache_nodes      = 1
  parameter_group_name = "default.redis3.2"
  engine_version       = "3.2.10"
  port                 = 6379

  log_delivery_configuration {
    destination      = aws_cloudwatch_log_group.example.name
    destination_type = "cloudwatch-logs"
    log_format       = "text"
    log_type         = "slow-log"
  }
  log_delivery_configuration {
    destination      = aws_kinesis_firehose_delivery_stream.redis.name
    destination_type = "kinesis-firehose"
    log_format       = "json"
    log_type         = "engine-log"
  }

  tags = {
    Name         = "${local.prefix}-redis-cluster"
    Environment  = var.Environment
    Organization = var.Organization
  }
}

resource "aws_cloudwatch_log_group" "redis" {
  name              = "${local.prefix}-redis-cluster-logs"
  retention_in_days = 7
}

resource "aws_kinesis_firehose_delivery_stream" "redis" {
  name        = "${local.prefix}-redis-cluster-logs-hose"
  destination = "extended_s3"

  extended_s3_configuration {
    role_arn   = aws_iam_role.firehose.arn
    bucket_arn = aws_s3_bucket.logs.arn
    prefix     = "redis-cluster-logs/"
  }
}

resource "aws_iam_role" "firehose" {
  name               = "${local.prefix}-firehose-role"
  assume_role_policy = <<EOF
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Effect": "Allow",
      "Principal": {
        "Service": "firehose.amazonaws.com"
      },
      "Action": "sts:AssumeRole",
      "Condition": {
        "StringEquals": {
          "sts:ExternalId": "${local.prefix}"
        }
      }
    }
  ]
}
EOF
}

resource "s3_bucket" "logs" {
  bucket = "${local.prefix}-logs"
  acl    = "private"
  tags = {
    Name         = "${local.prefix}-logs"
    Environment  = var.Environment
    Organization = var.Organization
  }
}
