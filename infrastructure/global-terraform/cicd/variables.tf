variable "project_id" {
  type = string
  description = "GCP Project ID for CICD related infra"
}

variable "artifiact_registry_project_id" {
  type = string
}

variable "artifact_registry_name" {
  type = string
}

variable "artifact_registry_location" {
  type = string
}