terraform {
  backend "remote" {
    hostname = "app.terraform.io"
    organization = "chroju"

    workspaces {
      name = "sample"
    }
  }
  required_version = "> 0.12.0, <= 0.12.22"
}
