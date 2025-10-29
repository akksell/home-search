output "github_ci_account" {
  value = google_service_account.github_actions_ci.email
}

output "github_ci_account_id" {
  value = google_service_account.github_actions_ci.account_id
}

output "workload_identity_pool_name" {
  value = google_iam_workload_identity_pool.github_actions_ci.name
}

output "workload_identity_pool_provider_name" {
  value = google_iam_workload_identity_pool_provider.github_actions_ci_provider.name
}