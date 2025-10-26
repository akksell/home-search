variable "service_name" {
  type = string
}

variable "environment" {
  type = string
  default = "dev"

  validation {
    condition = contains(["dev", "stage", "prod"], var.environment)
    error_message = "Invalid environment. Must be one of 'dev', 'stage', or 'prod'"
  }
}

variable "project_id" {
  type = string
}
