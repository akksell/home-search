module "tf_state_backend" {
  source = "../tf-backend"

  project_id = var.project_id
  service_name = "homesearch-global-infra"
  environment = "prod"
}

terraform {
  required_providers {
    google = {
      source = "hashicorp/google"
      version = "7.7.0"
    }
  }

  # IDK what happened because all of my other
  # infra uses 1.12.1 but I'm getting an error in my nix shell
  # so... ¯\_(ツ)_/¯
  required_version = "1.12.0"
}

provider "google" {
  project = var.project_id
  region = var.region
}