package setuptest

import (
	"fmt"

	"github.com/Azure/terratest-terraform-fluent/ops"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Output returns an Operative for the given output name.
// This allows us to perform assertions on the output value,
// e.g. ...Output("foo").HasValue("bar")
//
// This function works best with strongly typed values,
// e.g. `bool`, `number`, `string`, `list`, `map`, etc.
// If you use this with type `any`, then you
// will be dealing with strings and you assertion options
// will be limited.
func (resp Response) Output(name string) ops.Operative {
	ref := fmt.Sprintf("output.%s", name)
	allouts := terraform.OutputAll(resp.t, resp.Options)
	out, ok := allouts[name]
	return ops.Operative{
		Reference: ref,
		Exist:     ok,
		Actual:    out,
	}
}
