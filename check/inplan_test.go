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

func TestPlannedResourcesAre(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		expected  []string
		planned   *terraform.PlanStruct
		wantError bool
	}{
		{
			name:     "empty expected and empty plan",
			expected: []string{},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{},
			},
			wantError: false,
		},
		{
			name:     "matching single resource",
			expected: []string{"azurerm_resource_group.test"},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test": {},
				},
			},
			wantError: false,
		},
		{
			name: "matching multiple resources",
			expected: []string{
				"azurerm_resource_group.test",
				"azurerm_virtual_network.test",
				"azurerm_subnet.test",
			},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test":  {},
					"azurerm_virtual_network.test": {},
					"azurerm_subnet.test":          {},
				},
			},
			wantError: false,
		},
		{
			name:     "expected resource not in plan",
			expected: []string{"azurerm_resource_group.test"},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{},
			},
			wantError: true,
		},
		{
			name:     "unexpected resource in plan",
			expected: []string{},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test": {},
				},
			},
			wantError: true,
		},
		{
			name: "multiple expected resources not in plan",
			expected: []string{
				"azurerm_resource_group.test1",
				"azurerm_resource_group.test2",
			},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{},
			},
			wantError: true,
		},
		{
			name:     "multiple unexpected resources in plan",
			expected: []string{},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test1": {},
					"azurerm_resource_group.test2": {},
				},
			},
			wantError: true,
		},
		{
			name: "partial match - some expected missing, some unexpected present",
			expected: []string{
				"azurerm_resource_group.test1",
				"azurerm_resource_group.test2",
			},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test1": {},
					"azurerm_resource_group.test3": {},
				},
			},
			wantError: true,
		},
		{
			name: "expected subset of plan",
			expected: []string{
				"azurerm_resource_group.test1",
			},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test1": {},
					"azurerm_resource_group.test2": {},
				},
			},
			wantError: true,
		},
		{
			name: "plan subset of expected",
			expected: []string{
				"azurerm_resource_group.test1",
				"azurerm_resource_group.test2",
			},
			planned: &terraform.PlanStruct{
				ResourceChangesMap: map[string]*tfjson.ResourceChange{
					"azurerm_resource_group.test1": {},
				},
			},
			wantError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			pt := PlanType{Plan: tt.planned}
			err := pt.PlannedResourcesAre(tt.expected...)

			if tt.wantError {
				assert.Error(t, err.AsError())
			} else {
				assert.NoError(t, err.AsError())
			}
		})
	}
}
