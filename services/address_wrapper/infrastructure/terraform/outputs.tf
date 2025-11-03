output "geocode_api_secret_id" {
  value = google_secret_manager_secret.geocode_api.secret_id
}

output "service_account_email" {
  value = google_service_account.default.email
}

output "address_store_name" {
  value = google_firestore_database.address_store.name
}

output "address_store_id" {
  value = google_firestore_database.address_store.id
}
