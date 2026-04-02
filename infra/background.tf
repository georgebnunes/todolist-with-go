terraform {
  backend "s3" {
    bucket         = "geos-todos"
    key            = "golang-todos/terraform.tfstate"
    region         = "us-east-1"
    encrypt        = true
  }
}