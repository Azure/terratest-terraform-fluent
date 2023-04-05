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
  filename = "test.txt"
}


resource "local_file" "test_simple_json" {
  content  = local.jsondata_simple
  filename = "test_json.txt"
}

resource "local_file" "test_array_json" {
  content  = local.jsondata_array
  filename = "test_json.txt"
}

locals {
  jsondata_simple = jsonencode({
    test = "test"
  })
  jsondata_array = jsonencode(
    [{test = "test"}, {test2 = "test2"}]
  )
}
