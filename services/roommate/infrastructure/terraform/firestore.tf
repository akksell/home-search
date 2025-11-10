resource "google_firestore_database" "roommate_store" {
  name = "${var.service_name}-roommate-store-${var.environment}"
  location_id = var.region
  project = var.project_id
  type = "FIRESTORE_NATIVE"
  delete_protection_state = "DELETE_PROTECTION_ENABLED"
  deletion_policy = "DELETE"

}

resource "google_project_iam_member" "roommate_service_datastore_manager" {
  project = var.project_id
  role = "roles/datastore.user"
  member = "serviceAccount:${google_service_account.app_service_account.email}"

  condition {
    title = "restrict_rw_access_to_roommate_store"
    description = "Only access the roommate firestore instance - used for service accounts for the roommate service"
    expression = "resource.name.startsWith('projects/${var.project_id}/databases/${google_firestore_database.roommate_store.name}')"
  }
}