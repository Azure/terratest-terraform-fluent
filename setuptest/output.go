package setuptest

import (
	"fmt"

	"github.com/Azure/terratest-terraform-fluent/ops"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func (resp Response) Output(name string) ops.Operative {
	ref := fmt.Sprintf("output.%s", name)
	allouts := terraform.OutputAll(resp.t, resp.Options)
	actual, exist := allouts[name]
	return ops.Operative{
		Reference: ref,
		Exist:     exist,
		Actual:    actual,
	}
}
