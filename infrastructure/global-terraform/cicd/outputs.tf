output "github_ci_account" {
  value = google_service_account.github_actions_ci.email
}

output "github_ci_account_id" {
  value = google_service_account.github_actions_ci.account_id
}
