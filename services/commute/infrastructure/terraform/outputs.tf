output "firestore_id" {
  value = google_firestore_database.point_of_interest_store.id
}

output "service_account_id" {
  value = google_service_account.app_service_account.id
}

output "service_account_email" {
  value = google_service_account.app_service_account.email
}