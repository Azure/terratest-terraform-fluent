package check

import (
	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// ThatType is a type which can be used for more fluent assertions for a given Resource
type ThatType struct {
	Plan         *terraform.PlanStruct
	ResourceName string
}

// That returns a type which can be used for more fluent assertions for a given resource.
func (p PlanType) That(resourceName string) ThatType {
	return ThatType{
		Plan:         p.Plan,
		ResourceName: resourceName,
	}
}

// Exists returns an *testError.Error if the resource does not exist in the plan
func (t ThatType) Exists() *testerror.Error {
	if _, ok := t.Plan.ResourcePlannedValuesMap[t.ResourceName]; !ok {
		return testerror.Newf(
			"%s: resource not found in plan",
			t.ResourceName,
		)
	}
	return nil
}

// DoesNotExist returns an *testerror.Error if the resource exists in the plan
func (t ThatType) DoesNotExist() *testerror.Error {
	if _, exists := t.Plan.ResourcePlannedValuesMap[t.ResourceName]; exists {
		return testerror.Newf(
			"%s: resource found in plan",
			t.ResourceName,
		)
	}
	return nil
}

// Key returns a type which can be used for more fluent assertions for a given Resource & Key combination
func (t ThatType) Key(key string) ThatTypeWithKey {
	return ThatTypeWithKey{
		Plan:         t.Plan,
		ResourceName: t.ResourceName,
		Key:          key,
	}
}
