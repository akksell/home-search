# TODO: create a federated workload identity instead of using a service account
# service account was just easier to create and the first thing that popped in my
# head. This is less secure though
resource "google_service_account" "github_actions_ci" {
  project = var.project_id
  account_id = "github-actions-ci"
  display_name = "Github Actions"
}

resource "google_artifact_registry_repository_iam_member" "ci_artifact_repository_member" {
  project = var.artifiact_registry_project_id
  location = var.artifact_registry_location
  repository = var.artifact_registry_name
  role = "roles/artifactregistry.writer"
  member = "serviceAccount:${google_service_account.github_actions_ci.email}"
}
