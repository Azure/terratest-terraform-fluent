terraform {
  required_version = ">= 1.3.0"
}

variable "test" {
  type = string
  default = ""
}

output "test" {
  value = var.test
}
