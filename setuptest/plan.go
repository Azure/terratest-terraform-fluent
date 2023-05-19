package setuptest

import (
	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Apply runs terraform apply for the given Response and returns the error.
func (resp Response) Plan() *testerror.Error {
	_, err := terraform.PlanE(resp.t, resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// // Apply runs terraform apply for the given Response and returns the error.
// func (resp Response) Plan() *testerror.Error {
// 	_, err := terraform.Sho(resp.t, resp.Options)
// 	if err != nil {
// 		return testerror.New(err.Error())
// 	}
// 	return nil
// }
