variable "project_id" {
  description = "Project ID to contain the service"
  type = string
}

variable "region" {
  description = "GCP region to host and contain the service"
  type = string
  default = "us-central1"
}

variable "environment" {
  description = "Environment to deploy (dev, stage, prod)"
  type = string
  default = "dev"
}

variable "service_name" {
  description = "Name of the service"
  type = string
}
