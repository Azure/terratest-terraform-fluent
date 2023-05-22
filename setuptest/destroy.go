package setuptest

import (
	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Destroy runs terraform destroy for the given Response and returns the error.
func (resp Response) Destroy() *testerror.Error {
	_, err := terraform.DestroyE(resp.t, resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// DestroyWithRetry will retry the terraform destroy command up to the specified number of times.
func (resp Response) DestroyRetry(r Retry) *testerror.Error {
	resp.Options.RetryableTerraformErrors = map[string]string{
		".*": "Retry destroy on any error",
	}
	resp.Options.MaxRetries = r.Max
	resp.Options.TimeBetweenRetries = r.Wait
	_, err := terraform.DestroyE(resp.t, resp.Options)

	if err != nil {
		return testerror.Newf("terraform destroy failed after %d attempts: %v", r.Max, err)
	}
	return nil
}
