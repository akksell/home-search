resource "google_firestore_database" "point_of_interest_store" {
  name = "${var.service_name}-poi-store"
  location_id = var.region
  type = "FIRESTORE_NATIVE"
  delete_protection_state = "DELETE_PROTECTION_ENABLED"
  deletion_policy = "DELETE"

}

resource "google_project_iam_member" "compute_service_datastore_manager" {
  project = var.project_id
  role = "roles/datastore.user"
  member = "serviceAccount:${google_service_account.app_service_account.email}"

  condition {
    title = "restrict_rw_access_to_poi_store"
    description = "Only access the point_of_interest firestore instance - used for service accounts for the compute service"
    expression = "resource.name == '${google_firestore_database.point_of_interest_store.name}'"
  }
}