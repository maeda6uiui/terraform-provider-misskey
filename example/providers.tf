terraform {
  required_providers {
    misskey = {
      source  = "maeda6uiui/misskey"
      version = "0.0.1-alpha1"
    }
  }

  backend "local" {
    path = "terraform.tfstate"
  }

  required_version = "~>1.9"
}

provider "misskey" {
  server_url = "https://misskey-dabansky.com"
}
