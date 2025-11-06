resource "google_cloud_run_v2_service" "default" {
  name = "hs-${var.service_name}-${var.environment}"
  location = var.region
  description = "Service to get google place ids from an address"
  launch_stage = "ALPHA"

  labels = {
    component: "core"
  }

  ingress = "INGRESS_TRAFFIC_INTERNAL_ONLY"

  scaling {
    min_instance_count = 0
    max_instance_count = 2
    scaling_mode = "AUTOMATIC"
  }

  template {
    containers {
      image = "us-south1-docker.pkg.dev/home-search-475404/homesearch-services-docker/address_wrapper:ee4379a8"
      ports {
        name           = "h2c"
        container_port = 8080
      }
    }

    timeout = "10s"
    health_check_disabled = false

    service_account = google_service_account.default.email
  }

}