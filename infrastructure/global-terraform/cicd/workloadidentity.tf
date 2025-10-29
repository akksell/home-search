resource "google_iam_workload_identity_pool" "github_actions_ci" {
  project = var.project_id
  workload_identity_pool_id = "github-actions"
  display_name = "Github Actions CICD Pool"
}

resource "google_iam_workload_identity_pool_provider" "github_actions_ci_provider" {
  project = var.project_id
  workload_identity_pool_id = google_iam_workload_identity_pool.github_actions_ci.workload_identity_pool_id
  workload_identity_pool_provider_id = "github-actions-provider"
  description = "Github Actions CICD Identity Pool Provider"

  attribute_condition = <<EOT
      attribute.repository_owner_id == "5962960" &&
      attribute.repository_id == "1077226486" &&
      attribute.ref == "refs/heads/main" &&
      attribute.type == "branch"
  EOT

  attribute_mapping = {
    "google.subject"       = "assertion.sub"
    "attribute.actor"      = "assertion.actor"
    "attribute.aud"        = "assertion.aud"
    "attribute.repository_id" = "assertion.repository_id"
    "attribute.repository_owner_id" = "assertion.repository_owner_id"
    "attribute.ref" = "assertion.ref"
    "attribute.type" = "assertion.ref_type"
  }

  oidc {
    issuer_uri = "https://token.actions.githubusercontent.com"
  }
}

resource "google_artifact_registry_repository_iam_member" "ci_artifact_repository_member" {
  project = var.artifiact_registry_project_id
  location = var.artifact_registry_location
  repository = var.artifact_registry_name
  role = "roles/artifactregistry.writer"
  member = "principalSet://iam.googleapis.com/${google_iam_workload_identity_pool.github_actions_ci.name}/attribute.repository_id/1077226486"
}