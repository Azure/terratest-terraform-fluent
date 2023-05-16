package check

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestNumberOfResourcesInPlan(t *testing.T) {
	t.Parallel()

	pt := mockPlanType()
	err := pt.NumberOfResourcesEquals(2).AsError()
	assert.NoError(t, err)
	err = pt.NumberOfResourcesEquals(1).AsError()
	assert.ErrorContains(t, err, "expected 1 resources, got")
}

func TestInPlan(t *testing.T) {
	t.Parallel()
	ps := mockPlanStruct()
	ip := InPlan(ps)
	assert.Equal(t, ps, ip.Plan)
}

func TestThat(t *testing.T) {
	t.Parallel()

	mock := mockPlanType()
	t.Run("Exists", func(t *testing.T) {
		t.Parallel()
		tt := mock.That("test_resource")
		assert.True(t, tt.exists())
	})

	t.Run("NotExists", func(t *testing.T) {
		t.Parallel()
		tt := mock.That("not_exist")
		assert.False(t, tt.exists())
	})
}

func mockPlanStruct() *terraform.PlanStruct {
	return &terraform.PlanStruct{
		ResourcePlannedValuesMap: map[string]*tfjson.StateResource{
			"test_resource":  {},
			"test_resource2": {},
		},
	}
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
