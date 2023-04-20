package check

import (
	"encoding/json"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

// ThatTypeWithKey is a type which can be used for more fluent assertions for a given Resource & Key combination,
// together with a gjson query https://github.com/tidwall/gjson
type ThatTypeWithKeyQuery struct {
	Plan         *terraform.PlanStruct
	ResourceName string
	Key          string
	Query        string
}

// HasValue executes the provided gjson query on the resource and key combination
// and tests the result against the provided value.
// https://github.com/tidwall/gjson
func (twkq ThatTypeWithKeyQuery) HasValue(expected interface{}) *testerror.Error {
	resource := twkq.Plan.ResourcePlannedValuesMap[twkq.ResourceName]
	actual := resource.AttributeValues[twkq.Key]
	bytes, _ := json.Marshal(actual)
	result := gjson.GetBytes(bytes, twkq.Query)

	if err := validateEqualArgs(expected, result.Value()); err != nil {
		return testerror.Newf("invalid operation: %#v == %#v (%s)",
			expected,
			actual,
			err,
		)
	}

	if !assert.ObjectsAreEqualValues(result.Value(), expected) {
		return testerror.Newf(
			"%s: query result %s, for key %s not equal to assertion %s",
			twkq.ResourceName,
			result.Value(),
			twkq.Key,
			expected,
		)
	}
	return nil
}
