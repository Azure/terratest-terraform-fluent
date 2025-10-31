package check

import (
	"strings"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
)

// InPlan is the entry point for checking the Terraform plan.
func InPlan(plan *terraform.PlanStruct) PlanType {
	return PlanType{
		Plan: plan,
	}
}

// PlanType is a type which can be used for more fluent assertions on the Terraform plan.
type PlanType struct {
	Plan *terraform.PlanStruct
}

// NumberOfResourcesEquals checks that the number of resources in the plan is equal to the expected number.
func (p PlanType) NumberOfResourcesEquals(expected int) *testerror.Error {
	actual := len(p.Plan.ResourcePlannedValuesMap)
	if actual != expected {
		return testerror.Newf("expected %d resources, got %d", expected, actual)
	}
	return nil
}

// That returns a ThatType which can be used for more fluent assertions for a given resource.
func (p PlanType) That(resourceName string) ThatType {
	t := ThatType{
		Plan:         p.Plan,
		ResourceName: resourceName,
	}
	t.exists()
	return t
}

// PlannedResourcesAre takes a list of resource names and checks that they all exist in the plan.
func (p PlanType) PlannedResourcesAre(resourceNames ...string) *testerror.Error {
	found := make(map[string]struct{})
	msgs := make([]string, 0, len(resourceNames)+len(p.Plan.ResourceChangesMap))

	for _, resourceName := range resourceNames {
		found[resourceName] = struct{}{}
		if _, ok := p.Plan.ResourceChangesMap[resourceName]; !ok {
			msgs = append(msgs, "expected resource not found in plan: "+resourceName)
		}
	}

	for resourceName := range p.Plan.ResourceChangesMap {
		if _, ok := found[resourceName]; !ok {
			msgs = append(msgs, "unexpected resource found in plan: "+resourceName)
		}
	}

	if len(msgs) > 0 {
		return testerror.Newf("planned resources do not match expected:\n\n%s", strings.Join(msgs, "\n"))
	}

	return nil
}
