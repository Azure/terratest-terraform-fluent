terraform {
  required_version = ">= 1.3.0"
  required_providers {
    null = {
      source = "hashicorp/null"
      version = ">= 3.2.1"
    }
  }
}

resource "null_resource" "test" {
  provisioner "local-exec" {
    when = create
    command = "if [ -f 'ok' ]; then rm ok; fi"
    on_failure = fail
  }

  provisioner "local-exec" {
    when = destroy
    command = "if [ -f \"ok\" ]; then exit 0; else touch ok &&  exit 1; fi"
    on_failure = fail
  }
}
