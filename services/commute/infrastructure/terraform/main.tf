resource "google_cloud_run_v2_service" "default" {
  name = "hs-${var.service_name}-${var.environment}"
  location = var.region
  description = "Service to manage roommate details including roommate groups, points of interest, etc."
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
      name = "service"
      image = "us-south1-docker.pkg.dev/home-search-475404/homesearch-services-docker/commute:6509d955"
      /*
      ports {
        name           = "h2c"
        container_port = 8080
      }
      */

      env {
        name = "GOOGLE_PROJECT_ID"
        value = var.project_id
      }
      env {
        name = "COMMUTE_STORE_DATABASE"
        value = google_firestore_database.point_of_interest_store.name
      }
      env {
        name = "ADDRESS_WRAPPER_SERVICE_HOST"
        value = var.address_wrapper_service_host
      }
      env {
        name = "ROOMMATE_SERVICE_HOST"
        value = var.roommate_service_host
      }
    }

    timeout = "10s"
    health_check_disabled = false

    service_account = google_service_account.app_service_account.email
  }

}

resource "google_cloud_run_v2_service_iam_member" "address_service_member" {
  project = data.google_cloud_run_v2_service.address_service.project
  location = data.google_cloud_run_v2_service.address_service.location
  name = data.google_cloud_run_v2_service.address_service.name
  role = "roles/run.invoker"
  member = "serviceAccount:${google_service_account.app_service_account.email}"
}

data "google_cloud_run_v2_service" "address_service" {
  name = "hs-address-wrapper-${var.environment}"
  location = var.region
  project = var.project_id
}

resource "google_cloud_run_v2_service_iam_member" "roommate_service_member" {
  project = data.google_cloud_run_v2_service.roommate_service.project
  location = data.google_cloud_run_v2_service.roommate_service.location
  name = data.google_cloud_run_v2_service.roommate_service.name
  role = "roles/run.invoker"
  member = "serviceAccount:${google_service_account.app_service_account.email}"
}

data "google_cloud_run_v2_service" "roommate_service" {
  name = "hs-roommate-${var.environment}"
  location = var.region
  project = var.project_id
}
