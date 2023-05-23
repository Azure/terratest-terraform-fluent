terraform {
  required_version = ">= 1.4.0"
}


# Maps
variable "test_map" {
  type    = map(string)
  default = {}
}

output "test_map" {
  value = var.test_map
}

# String
variable "test_string" {
  type    = string
  default = ""
}

output "test_string" {
  value = var.test_string
}

# List
variable "test_list" {
  type    = list(string)
  default = []
}

output "test_list" {
  value = var.test_list
}

# Set
variable "test_set" {
  type    = set(string)
  default = []
}

output "test_set" {
  value = var.test_set
}

# Number
variable "test_number" {
  type    = number
  default = 0
}

output "test_number" {
  value = var.test_number
}

# Bool
variable "test_bool" {
  type    = bool
  default = false
}

output "test_bool" {
  value = var.test_bool
}
