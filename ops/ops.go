package ops

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"github.com/Azure/terratest-terraform-fluent/testerror"
	"github.com/stretchr/testify/assert"
	"github.com/tidwall/gjson"
)

// JsonAssertionFunc is a function which can be used to unmarshal a raw JSON message and check its contents.
type JsonAssertionFunc func(input json.RawMessage) (*bool, error)

// Operative is a type which can be used to compare the expected and actual values of a given combination
type Operative struct {
	Reference string
	Actual    any
	Exist     bool
	err       *testerror.Error
}

// Exists returns a non-nil *testerror.Error if the resource does not exist in the plan or if the key does not exist in the resource
func (o Operative) Exists() *testerror.Error {
	if err := isErrorOrNotExist(o); err != nil {
		return err
	}
	return nil
}

// DoesNotExist returns a non-nil *testerror.Error if the resource does not exist in the plan or if the key exists in the resource
func (o Operative) DoesNotExist() *testerror.Error {
	if o.Exist {
		return testerror.Newf(
			"%s: found when not expected",
			o.Reference,
		)
	}
	return nil
}

// HasValue returns a non-nil *testerror.Error if the resource does not exist in the plan
// or if the value of the key does not match the expected value
func (o Operative) HasValue(expected any) *testerror.Error {
	if err := isErrorOrNotExist(o); err != nil {
		return err
	}

	if err := validateEqualArgs(expected, o.Actual); err != nil {
		return testerror.Newf("invalid operation: %#v == %#v (%s)",
			expected,
			o.Actual,
			err,
		)
	}

	if !assert.ObjectsAreEqualValues(expected, o.Actual) {
		return testerror.Newf(
			"%s: expected value %v not equal to actual %v",
			o.Reference,
			expected,
			o.Actual,
		)
	}
	return nil
}

// ContainsString returns a non-nil *testerror.Error if the resource does not exist in the plan or if
// the value of the key does not contain the expected string
func (o Operative) ContainsString(expected string) *testerror.Error {
	if err := isErrorOrNotExist(o); err != nil {
		return err
	}

	actualString, ok := o.Actual.(string)
	if !ok {
		return testerror.Newf("Cannot convert value to string: %s", o.Reference)
	}

	if !strings.Contains(actualString, expected) {
		return testerror.Newf(
			"%s: expected value %s not contained within %s",
			o.Reference,
			expected,
			actualString,
		)
	}
	return nil
}

// GetValue returns the actual value and a *testerror.Error
func (o Operative) GetValue() (any, error) {
	if err := isErrorOrNotExist(o); err != nil {
		return nil, err
	}
	return o.Actual, nil
}

// ContainsJsonValue returns a *testerror.Error which asserts upon a given JSON string set
// by deserializing it and then asserting on it via the JsonAssertionFunc.
func (o Operative) ContainsJsonValue(assertion JsonAssertionFunc) *testerror.Error {
	if err := isErrorOrNotExist(o); err != nil {
		return err
	}

	if o.Actual == nil || o.Actual == "" {
		return testerror.Newf(
			"%s: is empty",
			o.Reference,
		)
	}

	actual, actualok := o.Actual.(string)
	if !actualok {
		return testerror.Newf(
			"%s: value is not a string",
			o.Reference,
		)
	}

	j := json.RawMessage(actual)
	assertok, err := assertion(j)
	if err != nil {
		return testerror.Newf(
			"%s: asserting value for %q: %+v",
			o.Reference,
			o.Actual,
			err,
		)
	}

	if assertok == nil || !*assertok {
		return testerror.Newf(
			"%s: assertion failed for %q",
			o.Reference,
			o.Actual,
		)
	}

	return nil
}

// Query executes the provided gjson query on the data in the actual value
// and overwrites the actual value with the result of the query.
// https://github.com/tidwall/gjson
func (o Operative) Query(query string) Operative {
	if err := isErrorOrNotExist(o); err != nil {
		o.Exist = false
		o.err = err
		return o
	}
	o.Reference = fmt.Sprintf("%s?%s", o.Reference, query)
	var bytes []byte
	// If the actual value is a string, we assume it is JSON and try to parse it.
	// Otherwise, we marshal it to JSON and try to parse it.
	actual, ok := o.Actual.(string)
	if ok {
		bytes = []byte(actual)
	} else {
		bytes, _ = json.Marshal(o.Actual)
	}

	if !gjson.ValidBytes(bytes) {
		o.err = testerror.Newf(
			"%s: actual value %s not valid JSON",
			o.Reference,
			o.Actual,
		)
		o.Actual = nil
		return o
	}

	o.Actual = gjson.GetBytes(bytes, query).Value()
	if o.Actual == nil {
		o.Exist = false
	}
	return o
}

// validateEqualArgs checks whether provided arguments can be safely used in the
// HasValue function.
func validateEqualArgs(expected, actual any) error {
	if expected == nil && actual == nil {
		return nil
	}

	if isFunction(expected) || isFunction(actual) {
		return fmt.Errorf("cannot take func type as argument")
	}
	return nil
}

func isFunction(arg any) bool {
	if arg == nil {
		return false
	}
	return reflect.TypeOf(arg).Kind() == reflect.Func
}

func isErrorOrNotExist(o Operative) *testerror.Error {
	if o.err != nil {
		return o.err
	}
	if !o.Exist {
		return testerror.Newf(
			"%s: not found when expected",
			o.Reference,
		)
	}
	return nil
}
