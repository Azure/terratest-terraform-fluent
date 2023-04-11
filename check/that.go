package check

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// ThatType is a type which can be used for more fluent assertions for a given Resource
type ThatType struct {
	Plan         *terraform.PlanStruct
	ResourceName string
}

// JsonAssertionFunc is a function which can be used to unmarshal a raw JSON message and check its contents.
type JsonAssertionFunc func(input json.RawMessage) (*bool, error)

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

// ThatTypeWithKey is a type which can be used for more fluent assertions for a given Resource & Key combination
type ThatTypeWithKey struct {
	Plan         *terraform.PlanStruct
	ResourceName string
	Key          string
}

// HasValue returns a *testerror.Error if the resource does not exist in the plan or if the value of the key does not match the
// expected value
func (twk ThatTypeWithKey) HasValue(expected interface{}) *testerror.Error {
	if err := twk.Exists(); err != nil {
		return err
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	actual := resource.AttributeValues[twk.Key]

	if err := validateEqualArgs(expected, actual); err != nil {
		return testerror.Newf("invalid operation: %#v == %#v (%s)",
			expected,
			actual,
			err,
		)
	}

	if !assert.ObjectsAreEqualValues(actual, expected) {
		return testerror.Newf(
			"%s: attribute %s, planned value %s not equal to assertion %s",
			twk.ResourceName,
			twk.Key,
			actual,
			expected,
		)
	}
	return nil
}

// ContainsJsonValue returns a *testerror.Error which asserts upon a given JSON string set into
// the State by deserializing it and then asserting on it via the JsonAssertionFunc.
func (twk ThatTypeWithKey) ContainsJsonValue(assertion JsonAssertionFunc) *testerror.Error {
	if err := twk.Exists(); err != nil {
		return err
	}

	if twk.HasValue("") == nil {
		return testerror.Newf(
			"%s: key %s was empty",
			twk.ResourceName,
			twk.Key,
		)
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	actual := resource.AttributeValues[twk.Key]
	j := json.RawMessage(actual.(string))
	ok, err := assertion(j)
	if err != nil {
		return testerror.Newf(
			"%s: asserting value for %q: %+v",
			twk.ResourceName,
			twk.Key,
			err,
		)
	}

	if ok == nil || !*ok {
		return testerror.Newf(
			"%s: assertion failed for %q: %+v",
			twk.ResourceName,
			twk.Key,
			err,
		)
	}

	return nil
}

// Exists returns a *testerror.Error if the resource does not exist in the plan or if the key does not exist in the resource
func (twk ThatTypeWithKey) Exists() *testerror.Error {
	if err := InPlan(twk.Plan).That(twk.ResourceName).Exists(); err != nil {
		return testerror.New(err.Error())
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	if _, exists := resource.AttributeValues[twk.Key]; !exists {
		return testerror.Newf(
			"%s: key %s not found in resource",
			twk.ResourceName,
			twk.Key,
		)
	}
	return nil
}

// DoesNotExist returns a *testerror.Error if the resource does not exist in the plan or if the key exists in the resource
func (twk ThatTypeWithKey) DoesNotExist() *testerror.Error {
	if err := InPlan(twk.Plan).That(twk.ResourceName).Exists(); err != nil {
		return testerror.Newf(err.Error())
	}

	resource := twk.Plan.ResourcePlannedValuesMap[twk.ResourceName]
	if _, exists := resource.AttributeValues[twk.Key]; exists {
		return testerror.Newf(
			"%s: key %s found in resource",
			twk.ResourceName,
			twk.Key,
		)
	}
	return nil
}

// validateEqualArgs checks whether provided arguments can be safely used in the
// HasValue function.
func validateEqualArgs(expected, actual interface{}) error {
	if expected == nil && actual == nil {
		return nil
	}

	if isFunction(expected) || isFunction(actual) {
		return fmt.Errorf("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg interface{}) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}
