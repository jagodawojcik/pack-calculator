terraform {
  required_version = ">= 1.14"

  required_providers {
    google = {
      source  = "hashicorp/google"
      version = ">= 5.0"
    }
  }

  backend "gcs" {
    bucket = "dev-terraform-state-61283"
  }
}

provider "google" {
  project                     = var.project_id
  region                      = var.region
  impersonate_service_account = var.iac_service_account_email
}


###################################
# Enable APIs
###################################
resource "google_project_service" "cloud_resource_manager_api" {
  project = var.project_id
  service = "cloudresourcemanager.googleapis.com"
}

resource "google_project_service" "artifact_registry_api" {
  project = var.project_id
  service = "artifactregistry.googleapis.com"
}

resource "google_project_service" "cloud_run_api" {
  project = var.project_id
  service = "run.googleapis.com"
}

###################################
# Google Artifact Registry
###################################
resource "google_artifact_registry_repository" "service_docker_repo" {
  project       = var.project_id
  location      = var.region
  repository_id = "packs-service-repo"
  description   = "Docker images for packs-service"
  format        = "DOCKER"
}

resource "google_project_iam_member" "artifact_registry_reader" {
  project = var.project_id
  role    = "roles/artifactregistry.reader"
  member  = "serviceAccount:${google_service_account.packs_server_runner_sa.email}"
}


#################################### 
# Cloud Run Service
####################################
resource "google_cloud_run_v2_service" "packs_server" {
  name     = "packs-service"
  location = var.region
  ingress  = "INGRESS_TRAFFIC_ALL"

  template {
    service_account = google_service_account.packs_server_runner_sa.email
    containers {
      image = "${var.region}-docker.pkg.dev/${var.project_id}/${google_artifact_registry_repository.service_docker_repo.name}/packs-server:dev"

      env {
        name  = "PACK_SIZES"
        value = join(",", var.pack_sizes)
      }
        resources {
        limits = {
            cpu    = "2"
            memory = "8Gi"
        }
        }
    }
  }
  scaling {
    max_instance_count = 10
  }
  deletion_protection = false # Keep experimental
}



resource "google_service_account" "packs_server_runner_sa" {
  account_id   = "packs-server-runner-sa"
  display_name = "Runner Service Account for Packs Service"
}

resource "google_cloud_run_v2_service_iam_member" "packs_service_invoker" {
  name = google_cloud_run_v2_service.packs_server.name
  project  = google_cloud_run_v2_service.packs_server.project
  location = google_cloud_run_v2_service.packs_server.location
  role     = "roles/run.invoker"
  member   = "allUsers"
}
