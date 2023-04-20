terraform {
  required_version = ">= 1.4.0"
}

resource "terraform_data" "test_map" {
  input = tomap({
    "test_key" = "test",
    "test_key_2" = "test2"
  })
}

resource "terraform_data" "test_map_list" {
  input = tomap({
    "test_key" = ["test", "test2"],
    "test_key_2" = ["testA", "testB"]
  })
}

resource "terraform_data" "test_nested_map" {
  input = tomap({
    "test_key" = tomap({
      "nested_key" = "test_nested"
    }),
    "test_key_2" = tomap({
      "nested_key2" = "test_nested2"
    }),
  })
}

resource "terraform_data" "invalid_json" {
  input = null
}
