resource "google_secret_manager_secret" "geocode_api" {
  project   = var.project_id
  secret_id = "${var.service_name}-geocode-api-key-${var.environment}"

  replication {
    auto {}
  }

  labels = {
    environment = var.environment
    owner       = var.service_name
  }

  annotations = {
    purpose = "application-api-key"
  }
}

resource "google_secret_manager_secret_iam_member" "geocode_api_member" {
  project   = google_secret_manager_secret.geocode_api.project
  secret_id = google_secret_manager_secret.geocode_api.secret_id
  role      = "roles/secretmanager.secretAccessor"
  member    = "serviceAccount:${google_service_account.default.email}"
}
