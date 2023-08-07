terraform {
  backend "gcs" {
    bucket = "go-wayback-infrastructure"
    prefix = "terraform/state"
  }
}
