package setuptest

import (
	"testing"
	"time"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"gopkg.in/matryer/try.v1"
)

// Destroy runs terraform destroy for the given Response and returns the error.
func (resp Response) Destroy(t *testing.T) *testerror.Error {
	_, err := terraform.DestroyE(t, resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// DestroyWithRetry will retry the terraform destroy command up to the specified number of times.
func (resp Response) DestroyRetry(t *testing.T, r Retry) *testerror.Error {
	if try.MaxRetries < r.Max {
		try.MaxRetries = r.Max
	}
	err := try.Do(func(attempt int) (bool, error) {
		_, err := terraform.DestroyE(t, resp.Options)
		if err != nil {
			t.Logf("terraform destroy failed attempt %d/%d: waiting %s", attempt, r.Max, r.Wait)
			if attempt < r.Max {
				time.Sleep(r.Wait)
			}
		}
		return attempt < r.Max, err
	})
	if err != nil {
		return testerror.Newf("terraform destroy failed after %d attempts: %v", r.Max, err)
	}
	return nil
}
