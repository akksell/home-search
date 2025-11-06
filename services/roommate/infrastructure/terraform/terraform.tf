module "tf_state_backend" {
  source = "../../../../infrastructure/tf-backend"
  
  service_name = var.service_name
  environment = var.environment
  project_id = var.project_id
}

terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "7.7.0"
    }
  }

  required_version = "1.12.0"
}

provider "google" {
  project = var.project_id
  region = var.region
}