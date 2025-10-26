resource "google_service_account" "app_service_account" {
  account_id = "roommate-${var.service_name}-${var.environment}"
  display_name = "Roommate API Service Account"
  description = "Manages/accesses resources related to Roommate service (e.g. firestore)"
  project = var.project_id
}