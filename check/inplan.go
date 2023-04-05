package check

import (
	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

func InPlan(plan *terraform.PlanStruct) PlanType {
	return PlanType{
		Plan: plan,
	}
}

type PlanType struct {
	Plan *terraform.PlanStruct
}

func (p PlanType) NumberOfResourcesEquals(expected int) *testerror.Error {
	actual := len(p.Plan.ResourcePlannedValuesMap)
	if actual != expected {
		return testerror.Newf("expected %d resources, got %d", expected, actual)
	}
	return nil
}
