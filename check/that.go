package check

import (
	"fmt"

	"github.com/Azure/terratest-terraform-fluent/ops"
	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// ThatType is a type which can be used for more fluent assertions for a given Resource
type ThatType struct {
	Plan         *terraform.PlanStruct
	ResourceName string
}

// Exists returns a *testError.Error if the resource does not exist in the plan
func (t ThatType) Exists() *testerror.Error {
	if !t.exists() {
		return testerror.Newf(
			"%s: resource not found in plan",
			t.ResourceName,
		)
	}
	return nil
}

func (t *ThatType) exists() bool {
	if _, ok := t.Plan.ResourcePlannedValuesMap[t.ResourceName]; !ok {
		return false
	}
	return true
}

// DoesNotExist returns an *testerror.Error if the resource exists in the plan
func (t ThatType) DoesNotExist() *testerror.Error {
	if t.exists() {
		return testerror.Newf(
			"%s: resource found in plan",
			t.ResourceName,
		)
	}
	return nil
}

// Key returns an ops.Operative type which can be used to compare and query the data
func (t ThatType) Key(key string) ops.Operative {
	ref := fmt.Sprintf("%s.%s", t.ResourceName, key)

	if !t.exists() {
		return ops.Operative{
			Exist:     false,
			Reference: ref,
		}
	}

	actual, ok := t.Plan.ResourcePlannedValuesMap[t.ResourceName].AttributeValues[key]
	if !ok {
		return ops.Operative{
			Exist:     false,
			Reference: ref,
		}
	}

	return ops.Operative{
		Exist:     true,
		Reference: ref,
		Actual:    actual,
	}
}
