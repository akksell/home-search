variable "project_id" {
  type = string
  description = "GCP project id"
}

variable "region" {
  type = string 
  description = "GCP region to deploy registry to"
}

variable "repository_id" {
  type = string
  description = "Name of the Artifact Repository"
}