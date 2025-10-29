resource "google_artifact_registry_repository" "default" {
  location = var.region
  project = var.project_id
  repository_id = var.repository_id
  description = "Repository to store service container images"
  format = "DOCKER"
}