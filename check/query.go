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
func (twkq ThatTypeWithKeyQuery) HasValue(expected any) *testerror.Error {

	err := ThatTypeWithKey{
		Plan:         twkq.Plan,
		ResourceName: twkq.ResourceName,
		Key:          twkq.Key,
	}.Exists()

	if err != nil {
		return err
	}

	var bytes []byte

	resource := twkq.Plan.ResourcePlannedValuesMap[twkq.ResourceName]
	actual := resource.AttributeValues[twkq.Key]

	// If the actual value is a string, we assume it is JSON and try to parse it.
	// Otherwise, we marshal it to JSON and try to parse it.
	actualS, ok := actual.(string)
	if !ok {
		bytes, _ = json.Marshal(actual)
	} else {
		bytes = []byte(actualS)
	}

	if !gjson.ValidBytes(bytes) {
		return testerror.Newf(
			"%s: attribute %s, planned value %s not valid JSON",
			twkq.ResourceName,
			twkq.Key,
			actual,
		)
	}

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
			"%s: query result %v, for key %s not equal to assertion %v",
			twkq.ResourceName,
			result.Value(),
			twkq.Key,
			expected,
		)
	}
	return nil
}
