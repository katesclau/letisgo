variable "Environment" {
  description = "The environment for the infrastructure (e.g., dev, staging, prod)"
  type        = string
  default     = "local"
}

variable "ServiceName" {
  description = "The name of the service being deployed"
  type        = string
  default     = "backend"
}

variable "Organization" {
  description = "The organization name"
  type        = string
  default     = "letisgo"
}

locals {
  prefix = "${var.Organization}-${var.ServiceName}-${var.Environment}"
}
