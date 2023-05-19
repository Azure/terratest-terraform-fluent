terraform {
  required_version = ">= 1.4.0"
}

variable "test" {
}

output "test" {
  value = var.test
}
