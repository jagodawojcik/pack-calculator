variable "iac_service_account_email" {
  description = "Service account for Terraform deployment"
  type        = string
}

variable "project_id" {
  description = "GCP project ID to deploy resources into"
  type        = string
}

variable "region" {
  description = "GCP region to deploy resources into"
  type        = string
  default     = "europe-west4"
}

variable "pack_sizes" {
  description = "List of pack sizes"
  type        = list(string)
  default     = ["250", "500", "1000", "2000", "5000"]
}
