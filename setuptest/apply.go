package setuptest

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/retry"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// Retry is a configuration for retrying a terraform command.
// Max is the number of times to retry.
// Wait is the amount of time to wait between each retry.
type Retry struct {
	Max  int
	Wait time.Duration
}

// DefaultRetry is the default retry configuration.
// It will retry up to 5 times with a 1 minute wait between each attempt.
var DefaultRetry = Retry{
	Max:  5,
	Wait: time.Minute,
}

// DefaultRetry is the faster retry configuration.
// It will retry up to 6 times with a 20 second wait between each attempt.
var FastRetry = Retry{
	Max:  6,
	Wait: 20 * time.Second,
}

// DefaultRetry is the slower retry configuration.
// It will retry up to 15 times with a 2 minute wait between each attempt.
var SlowRetry = Retry{
	Max:  15,
	Wait: 2 * time.Minute,
}

// Apply runs terraform apply for the given Response and returns the error.
func (resp Response) Apply() *testerror.Error {
	// If there's no plan file, then we need to run apply without a plan.
	if _, err := os.Stat(filepath.Join(resp.Options.TerraformDir, resp.Options.PlanFilePath)); err != nil {
		if !os.IsNotExist(err) {
			return testerror.New(err.Error())
		}
		opts := resp.Options
		opts.PlanFilePath = ""
		_, err = terraform.ApplyE(resp.t, opts)
		if err != nil {
			return testerror.New(err.Error())
		}
		return nil
	}

	_, err := terraform.ApplyE(resp.t, resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// Apply runs terraform apply, then plan for the given Response and checks for any changes,
// it then returns the error.
func (resp Response) ApplyIdempotent() *testerror.Error {
	_, err := terraform.ApplyAndIdempotentE(resp.t, resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// Apply runs terraform apply, then performs a retry loop with a plan.
// If the configuration is not idempotent, it will retry up to the specified number of times.
// It then returns the error.
func (resp Response) ApplyIdempotentRetry(r Retry) *testerror.Error {
	_, err := terraform.ApplyE(resp.t, resp.Options)

	if err != nil {
		return testerror.New(err.Error())
	}

	_, err = retry.DoWithRetryE(resp.t, "terraform plan", r.Max, r.Wait, func() (string, error) {
		exitCode, err := terraform.PlanExitCodeE(resp.t, resp.Options)
		if err != nil {
			return "", retry.FatalError{Underlying: err}
		}
		if exitCode == 1 {
			return "", retry.FatalError{Underlying: errors.New("terraform plan exit code 1")}
		}
		if exitCode == 2 {
			return "", errors.New("terraform configuration not idempotent")
		}
		return "", nil
	})

	if err != nil {
		return testerror.New(err.Error())
	}

	return nil
}
