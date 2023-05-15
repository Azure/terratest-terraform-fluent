package check

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
)

func TestResourceExists(t *testing.T) {
	t.Parallel()

	tt := mockThatType()
	err := tt.Exists().AsError()
	assert.NoError(t, err)
}

func TestResourceExistsFail(t *testing.T) {
	t.Parallel()

	tt := mockThatType()
	tt.ResourceName = "not_exists"
	err := tt.Exists().AsError()
	assert.Error(t, err)
}

func TestResourceDoesNotExist(t *testing.T) {
	t.Parallel()

	tt := mockThatType()
	tt.ResourceName = "not_exists"
	err := tt.DoesNotExist().AsError()
	assert.NoError(t, err)
}

func TestResourceDoesNotExistFail(t *testing.T) {
	t.Parallel()

	tt := mockThatType()
	err := tt.DoesNotExist().AsError()
	assert.ErrorContains(t, err, "test_resource: resource found in plan")
}

func mockThatType() ThatType {
	return ThatType{
		Plan: &terraform.PlanStruct{
			ResourcePlannedValuesMap: map[string]*tfjson.StateResource{
				"test_resource": {
					AttributeValues: map[string]interface{}{},
				},
			},
		},
		ResourceName: "test_resource",
	}
}
