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

// SlowRetry is the slower retry configuration.
// It will retry up to 15 times with a 2 minute wait between each attempt.
var SlowRetry = Retry{
	Max:  15,
	Wait: 2 * time.Minute,
}

// Apply runs terraform apply for the given Response and returns the error.
// If the plan file does not exist, it will run terraform apply without a plan file.
func (resp Response) Apply() *testerror.Error {
	opts, err := checkPlanFileExists(resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	_, err = terraform.ApplyE(resp.t, opts)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// Apply runs terraform apply, then plan for the given Response and checks for any changes,
// it then returns the error.
// If the plan file does not exist, it will run terraform apply without a plan file.
func (resp Response) ApplyIdempotent() *testerror.Error {
	opts, err := checkPlanFileExists(resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}
	_, err = terraform.ApplyAndIdempotentE(resp.t, opts)
	if err != nil {
		return testerror.New(err.Error())
	}
	return nil
}

// Apply runs terraform apply, then performs a retry loop with a plan.
// If the configuration is not idempotent, it will retry up to the specified number of times.
// It then returns the error.
// If the plan file does not exist, it will run terraform apply without a plan file.
func (resp Response) ApplyIdempotentRetry(r Retry) *testerror.Error {
	opts, err := checkPlanFileExists(resp.Options)
	if err != nil {
		return testerror.New(err.Error())
	}

	_, err = terraform.ApplyE(resp.t, opts)
	if err != nil {
		return testerror.New(err.Error())
	}

	_, err = retry.DoWithRetryE(resp.t, "terraform plan", r.Max, r.Wait, func() (string, error) {
		exitCode, err := terraform.PlanExitCodeE(resp.t, opts)
		if err != nil {
			return "", retry.FatalError{Underlying: err}
		}
		switch exitCode {
		case 0:
			return "", nil
		case 2:
			return "", retry.FatalError{Underlying: errors.New("terraform configuraiton not idempotent")}
		default:
			return "", errors.New("terraform plan error")
		}
	})

	if err != nil {
		return testerror.New(err.Error())
	}

	return nil
}

// checkPlanFileExists takes in a terraform.Options and checks if the plan file exists.
// If it does not it returns a new terraform.Options with the PlanFilePath set to "" (to enable apply to be run without a plan file).
// If it does exist, it returns the original terraform.Options.
func checkPlanFileExists(opts *terraform.Options) (*terraform.Options, error) {
	if _, err := os.Stat(filepath.Join(opts.TerraformDir, opts.PlanFilePath)); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
		newopts := new(terraform.Options)
		*newopts = *opts
		newopts.PlanFilePath = ""
		return newopts, nil
	}
	return opts, nil
}
