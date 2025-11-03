resource "google_firestore_database" "address_store" { 
  name = "${var.service_name}-address-store-${var.environment}"
  location_id = var.region
  project = var.project_id
  type = "FIRESTORE_NATIVE"
  delete_protection_state = "DELETE_PROTECTION_ENABLED"
  deletion_policy = "DELETE"
}

resource "google_project_iam_member" "address_store_manager" {
  project = var.project_id
  role = "roles/datastore.user"
  member = "serviceAccount:${google_service_account.default.email}"

  condition {
    title = "restrict_rw_access_to_address_store"
    description = "Only access the address wrapper store firestore instance"
    expression = "resource.name == '${google_firestore_database.address_store.name}'"
  }
}
