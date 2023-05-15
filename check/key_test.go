package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/Azure/terratest-terraform-fluent/to"
	"github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHasValueInvalidArgs(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))
	err := twk.HasValue(func() {})
	assert.ErrorContains(t, err, "invalid operation")
}

func TestHasValueStrings(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("test")
	err := twk.HasValue("test")
	assert.NoError(t, err.AsError())
}

func TestHasValueStringsNotEqualError(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("test")
	err := twk.HasValue("not_equal")
	assert.ErrorContains(t, err, "attribute test_key, planned value test not equal to assertion not_equal")
}

func TestHasValueStringsToInt(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("123")
	err := twk.HasValue(123)
	assert.Error(t, err.AsError())
}

func TestKeyNotExistsError(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))
	twk.Key = "not_exists"
	err := twk.Exists()
	assert.ErrorContains(t, err, "key not_exists not found in resource")

}

func TestKeyNotExists(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))
	twk.Key = "not_exists"
	err := twk.DoesNotExist()
	assert.NoError(t, err.AsError())
}

func TestKeyNotExistsFail(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))
	err := twk.DoesNotExist()
	assert.ErrorContains(t, err, "key test_key found in resource")
}

func TestInSubdir(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs("testdata/test-in-subdir", "subdir").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("module.test.local_file.test").Key("content").HasValue("test").ErrorIsNil(t)
}

func TestInSubdirFail(t *testing.T) {
	t.Parallel()

	_, err := setuptest.Dirs("testdata/test-in-subdir", "not_exist").WithVars(nil).InitPlanShow(t)
	require.True(t, os.IsNotExist(err))
}

func TestJsonArrayAssertionFunc(t *testing.T) {
	t.Parallel()

	f := func(input json.RawMessage) (*bool, error) {
		i := make([]any, 0, 1)
		if err := json.Unmarshal(input, &i); err != nil {
			return nil, fmt.Errorf("JSON input is not an array")
		}
		if len(i) == 0 {
			return nil, fmt.Errorf("JSON input is empty")
		}
		if i[0].(map[string]any)["test"] != "test" {
			return nil, fmt.Errorf("JSON input key name is not equal to test")
		}

		return to.Ptr(true), nil
	}

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("[{\"test\":\"test\"}]")
	err := twk.ContainsJsonValue(JsonAssertionFunc(f))
	require.NoError(t, err.AsError())
}

func TestJsonEmpty(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return to.Ptr(true), nil
		},
	)

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("")
	twk.ContainsJsonValue(f).ErrorContains(t, "key test_key was empty")
}

func TestJsonAssertionFuncError(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return to.Ptr(true), errors.New("test error")
		},
	)

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("{\"test\": \"test\"}")
	twk.ContainsJsonValue(f).ErrorContains(t, "test error")
}

func TestJsonAssertionFuncFalse(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return to.Ptr(false), nil
		},
	)

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("{\"test\": \"test\"}")
	twk.ContainsJsonValue(f).ErrorContains(t, "assertion failed for \"test_key\"")
}

func TestJsonAssertionFuncNil(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return nil, nil
		},
	)

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("{\"test\": \"test\"}")
	err := twk.ContainsJsonValue(f)
	require.ErrorContains(t, err, "assertion failed for \"test_key\"")
}

func TestJsonSimpleAssertionFunc(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			i := make(map[string]any)
			if err := json.Unmarshal(input, &i); err != nil {
				return nil, fmt.Errorf("JSON input is not an map")
			}
			if len(i) == 0 {
				return nil, fmt.Errorf("JSON input is empty")
			}
			if i["test"] != "test" {
				return to.Ptr(false), nil
			}
			return to.Ptr(true), nil
		},
	)

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("{\"test\": \"test\"}")
	twk.ContainsJsonValue(f).ErrorIsNil(t)
}

func TestKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))
	twk.Key = "not_exist"
	err := twk.DoesNotExist()
	require.NoError(t, err.AsError())
}

func TestKeyDoesNotExistFail(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))

	err := twk.DoesNotExist()
	require.ErrorContains(t, err, "test_resource: key test_key found in resource")
}

func TestValidateEqualArgs(t *testing.T) {
	require.Nil(t, validateEqualArgs(nil, nil))
}

func TestValidateEqualArgsFuncFail(t *testing.T) {
	f1 := func() {}
	f2 := func() {}
	assert.ErrorContains(t, validateEqualArgs(f1, nil), "cannot take func type as argument")
	assert.ErrorContains(t, validateEqualArgs(nil, f2), "cannot take func type as argument")
}

func TestIsFunction(t *testing.T) {
	f := func() {}
	assert.True(t, isFunction(f))
}

func TestIsFunctionNot(t *testing.T) {
	i := 1
	assert.False(t, isFunction(i))
	var s *string
	assert.False(t, isFunction(s))
}

func TestContainsString(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("test")

	// Test that the function returns nil when the expected string is contained in the actual string
	err := twk.ContainsString("test")
	assert.Nil(t, err.AsError())

	// Test that the function returns an error when the expected string is not contained in the actual string
	err = twk.ContainsString("not in string")
	assert.NotNil(t, err.AsError())
	assert.Contains(t, err.Error(), "does not contain assertion")
}

func TestContainsStringNotAString(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey(any(nil))

	// Test that the function returns expected error if string conversion is not possible
	err := twk.ContainsString("test")
	assert.ErrorContains(t, err, "Cannot convert value to string")
}

func TestContainsStringKeyNotExists(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKey object
	twk := mockThatTypeWithKey("this is a test string")
	twk.Key = "not_exists"

	// Test that the function returns expected error if string conversion is not possible
	err := twk.ContainsString("test")
	assert.ErrorContains(t, err, "key not_exists not found in resource")
}

func mockThatTypeWithKey(val any) ThatTypeWithKey {
	return ThatTypeWithKey{
		Plan: &terraform.PlanStruct{
			ResourcePlannedValuesMap: map[string]*tfjson.StateResource{
				"test_resource": {
					AttributeValues: map[string]any{
						"test_key": val,
					},
				},
			},
		},
		ResourceName: "test_resource",
		Key:          "test_key",
	}
}
