variable "project_id" {
  description = "The ID of the Google project"
  default     = "go-wayback"
  type        = string
}

variable "region" {
  description = "The GCP region"
  default     = "europe-west3"
}
