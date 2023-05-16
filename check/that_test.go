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

func TestKey(t *testing.T) {
	t.Parallel()

	t.Run("Success", func(t *testing.T) {
		t.Parallel()
		tt := mockThatType()
		o := tt.Key("key")
		assert.Equal(t, "value", o.Actual)
		assert.True(t, o.Exist)
	})

	t.Run("ResourceNotFound", func(t *testing.T) {
		t.Parallel()
		tt := mockThatType()
		tt.ResourceName = "not_exists"
		o := tt.Key("key")
		assert.Nil(t, o.Actual)
		assert.False(t, o.Exist)
	})

	t.Run("KeyNotFound", func(t *testing.T) {
		t.Parallel()
		tt := mockThatType()
		o := tt.Key("not_exists")
		assert.Nil(t, o.Actual)
		assert.False(t, o.Exist)
	})
}

func mockThatType() ThatType {
	return ThatType{
		Plan: &terraform.PlanStruct{
			ResourcePlannedValuesMap: map[string]*tfjson.StateResource{
				"test_resource": {
					AttributeValues: map[string]any{
						"key": "value",
					},
				},
			},
		},
		ResourceName: "test_resource",
	}
}
