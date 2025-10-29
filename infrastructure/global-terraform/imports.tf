module "artifactregistry" {
  source = "./artifact-registry"
  project_id = var.project_id
  region = var.region
  repository_id = "homesearch-services-docker"
}

module "cicd" {
  source = "./cicd"

  project_id = var.project_id
  artifact_registry_location = module.artifactregistry.registry_location
  artifact_registry_name = module.artifactregistry.registry_name
  artifiact_registry_project_id = module.artifactregistry.registry_project
}