provider "google" {
  project = var.project_id
  region = var.region
}

/*resource "google_cloud_run_v2_service" "default" {
  name = var.service_name
  location = var.region
  deletion_protection = true
  launch_stage = "ALPHA"
  ingress = "INGRESS_TRAFFIC_ALL"
  scaling {
    min_instance_count = 0
    max_instance_count = 2
    scaling_mode = "AUTOMATIC"
  }
  template {
    containers {
      name = "${var.service_name}-app"
      image = # todo
      depends_on = ["envoy-auth-sidecar"]
    }
    containers {
      name = "envoy-auth-sidecar"
      image = # todo: envoy proxy
    }

    service_account = google_service_account.app_service_account.email
  }
}
*/