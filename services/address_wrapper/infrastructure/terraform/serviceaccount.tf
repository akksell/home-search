resource "google_service_account" "default" {
  account_id = "${var.service_name}-manager-${var.environment}"
  display_name = "Address Wrapper Account"
  description = "Manages/accesses resources required for address wrapper service to function"
  project = var.project_id
}
