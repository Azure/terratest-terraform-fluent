terraform {
  required_version = ">= 1.3.0"
  required_providers {
    local = {
      source  = "hashicorp/local"
      version = "2.3.0"
    }
  }
}
resource "local_file" "test" {
  content  = "test"
  filename = "test.txt"
}

resource "local_file" "test_int" {
  content  = 123
  filename = "test_int.txt"
}
