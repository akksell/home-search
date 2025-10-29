# NOTE: this is temporary in case I need it for an future CICD
#       interactions with GCP services - prefer to use workload
#       identity federation
resource "google_service_account" "github_actions_ci" {
  project = var.project_id
  account_id = "github-actions-ci"
  display_name = "Github Actions"
}

