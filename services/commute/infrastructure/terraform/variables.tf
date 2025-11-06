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

variable "address_wrapper_service_host" {
  description = "Internal GCP URL of the deployed Address Wrapper Service"
  type = string 
}

variable "roommate_service_host" {
  description = "Internal GCP URL of the deployed Roommate Service"
  type = string 
}
