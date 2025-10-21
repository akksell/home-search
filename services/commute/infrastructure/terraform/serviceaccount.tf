resource "google_service_account" "app_service_account" {
  account_id = "cr-${var.service_name}"
  display_name = "Commute API Service Account"
  description = "Manages resources related to commute service (e.g. routes api, firestore)"
  project = var.project_id
}
