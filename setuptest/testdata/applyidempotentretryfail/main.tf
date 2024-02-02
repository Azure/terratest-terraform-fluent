terraform {
  required_version = ">= 1.4.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = ">= 2.4.0"
    }
  }
}

resource "local_file" "test" {
  content         = "test"
  filename        = "test.txt"
  file_permission = "0644"
}

resource "terraform_data" "test" {
  lifecycle {
    replace_triggered_by = [
      local_file.test
    ]

  }
  provisioner "local-exec" {
    command = "rm -rf .terraform"
    when    = create
  }
  depends_on = [
    local_file.test,
  ]
}
