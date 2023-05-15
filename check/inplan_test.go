package check

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
)

func TestNumberOfResourcesInPlan(t *testing.T) {
	t.Parallel()

	pt := mockPlanType()
	pt.NumberOfResourcesEquals(2).ErrorIsNil(t)
}

func TestNumberOfResourcesInPlanWithError(t *testing.T) {
	t.Parallel()

	pt := mockPlanType()
	pt.NumberOfResourcesEquals(1).ErrorContains(t, "expected 1 resources, got")
}

func mockPlanType() PlanType {
	return PlanType{
		Plan: &terraform.PlanStruct{
			ResourcePlannedValuesMap: map[string]*tfjson.StateResource{
				"test_resource":  {},
				"test_resource2": {},
			},
		},
	}
}
