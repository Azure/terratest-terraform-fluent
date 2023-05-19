package setuptest

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"gopkg.in/matryer/try.v1"
)

// Retry is a configuration for retrying a terraform command.
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
// It will retry up to 20 times with a 1 minute wait between each attempt.
var SlowRetry = Retry{
	Max:  20,
	Wait: 1 * time.Minute,
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

	if try.MaxRetries < r.Max {
		try.MaxRetries = r.Max
	}
	err = try.Do(func(attempt int) (bool, error) {
		exitCode, err := terraform.PlanExitCodeE(resp.t, resp.Options)
		if err != nil {
			resp.t.Logf("terraform plan failed attempt %d/%d: waiting %s", attempt, r.Max, r.Wait)
			time.Sleep(r.Wait)
		}
		if exitCode == 2 {
			resp.t.Logf("terraform not idempotent attempt %d/%d: waiting %s", attempt, r.Max, r.Wait)
			err = errors.New("terraform configuration not idempotent")
			if attempt < r.Max {
				time.Sleep(r.Wait)
			}
		}

		return attempt < r.Max, err
	})

	if err != nil {
		return testerror.New(err.Error())
	}

	return nil
}
