resource "google_storage_bucket" "default" {
  name = "${var.service_name}-bucket-tf-state-${var.environment}"
  location = "US"
  storage_class = "STANDARD"
  force_destroy = false
  public_access_prevention = "enforced"
  project = var.project_id

  versioning {
    enabled = true
  }
}

resource "local_file" "default" {
  file_permission = "0644"
  filename        = "config/backend.${var.environment}.config"

  content = templatefile("${path.module}/backend.tftpl", { bucket_name = google_storage_bucket.default.name })
}