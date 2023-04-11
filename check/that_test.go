package check

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"testing"

	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	basicTestData = "testdata/basic"
)

func Bool(b bool) *bool {
	var b2 = b
	return &b2
}

func TestHasValueInvalidArgs(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue(func() {}).ErrorContains(t, "invalid operation")
}

func TestHasValueStrings(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue("test").ErrorIsNil(t)
}

func TestHasValueStringsNotEqualError(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("local_file.test").Key("content").HasValue("throwError"),
		"attribute content, planned value test not equal to assertion throwError",
	)
}

func TestHasValueStringsToInt(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.Error(
		t,
		InPlan(tftest.Plan).That("local_file.test_int").Key("content").HasValue(123).AsError(),
	)
}

func TestKeyNotExistsError(t *testing.T) {
	t.Parallel()

	tftest, _ := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("local_file.test").Key("not_exists").Exists(),
		"key not_exists not found in resource",
	)
}

func TestResourceExists(t *testing.T) {
	t.Parallel()

	tftest, _ := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	assert.NoErrorf(
		t,
		InPlan(tftest.Plan).That("local_file.test").Exists().AsError(),
		"resource local_file.test not found in plan",
	)
}

func TestResourceExistsFail(t *testing.T) {
	t.Parallel()

	tftest, _ := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	assert.Errorf(
		t,
		InPlan(tftest.Plan).That("not_exists").Exists(),
		"resource not_exists found in plan",
	)
}

func TestKeyNotExists(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	require.NoError(t, err)
	InPlan(tftest.Plan).That("local_file.test").Key("not_exists").DoesNotExist().ErrorIsNil(t)
}

func TestKeyNotExistsFail(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	defer tftest.Cleanup()
	require.NoError(t, err)
	require.Errorf(t, InPlan(tftest.Plan).That("local_file.test").Key("content").DoesNotExist(), "key content not found in resource when it should be")
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
		i := make([]interface{}, 0, 1)
		if err := json.Unmarshal(input, &i); err != nil {
			return nil, fmt.Errorf("JSON input is not an array")
		}
		if len(i) == 0 {
			return nil, fmt.Errorf("JSON input is empty")
		}
		if i[0].(map[string]interface{})["test"] != "test" {
			return nil, fmt.Errorf("JSON input key name is not equal to test")
		}

		return Bool(true), nil
	}

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_array_json").Key("content").ContainsJsonValue(JsonAssertionFunc(f)).ErrorIsNil(t)
}

func TestJsonEmpty(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return Bool(false), nil
		},
	)

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_empty_json").Key("content").ContainsJsonValue(f).ErrorContains(t, "key content was empty")
}

func TestJsonAssertionFuncError(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return Bool(false), errors.New("test error")
		},
	)

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_simple_json").Key("content").ContainsJsonValue(f).ErrorContains(t, "test error")
}

func TestJsonAssertionFuncFalse(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return Bool(false), nil
		},
	)

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_simple_json").Key("content").ContainsJsonValue(f).ErrorContains(t, "assertion failed for \"content\"")
}

func TestJsonAssertionFuncNil(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			return nil, nil
		},
	)

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_simple_json").Key("content").ContainsJsonValue(f).ErrorContains(t, "assertion failed for \"content\"")
}

func TestJsonSimpleAssertionFunc(t *testing.T) {
	t.Parallel()

	f := JsonAssertionFunc(
		func(input json.RawMessage) (*bool, error) {
			i := make(map[string]interface{})
			if err := json.Unmarshal(input, &i); err != nil {
				return nil, fmt.Errorf("JSON input is not an map")
			}
			if len(i) == 0 {
				return nil, fmt.Errorf("JSON input is empty")
			}
			if i["test"] != "test" {
				return Bool(false), nil
			}
			return Bool(true), nil
		},
	)

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test_simple_json").Key("content").ContainsJsonValue(f).ErrorIsNil(t)
}

func TestResourceDoesNotExist(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("not_exist").DoesNotExist().ErrorIsNil(t)
}

func TestResourceDoesNotExistFail(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	err = InPlan(tftest.Plan).That("local_file.test").DoesNotExist()
	require.ErrorContains(t, err, "local_file.test: resource found in plan")
}

func TestKeyDoesNotExist(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).That("local_file.test").Key("not_exist").DoesNotExist().ErrorIsNil(t)
}

func TestKeyDoesNotExistFail(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	err = InPlan(tftest.Plan).That("local_file.test").Key("content").DoesNotExist()
	require.ErrorContains(t, err, "local_file.test: key content found in resource")
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
