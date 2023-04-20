package check

import (
	"testing"

	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	queryTestData = "testdata/query"
)

func TestQueryMap(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.NoError(
		t,
		InPlan(tftest.Plan).That("terraform_data.test_map").Key("input").Query("test_key").HasValue("test").AsError(),
	)
}

func TestQueryMapLength(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.NoError(
		t,
		InPlan(tftest.Plan).That("terraform_data.test_map_list").Key("input").Query("test_key.#").HasValue(2).AsError(),
	)
}

func TestQueryNestedMap(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.NoError(
		t,
		InPlan(tftest.Plan).That("terraform_data.test_nested_map").Key("input").Query("test_key.nested_key").HasValue("test_nested").AsError(),
	)
}

func TestQueryNil(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.NoError(
		t,
		InPlan(tftest.Plan).That("terraform_data.invalid_json").Key("input").Query(".").HasValue(nil).AsError(),
	)
}

func TestQueryInvalidArgs(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()

	InPlan(tftest.Plan).That("terraform_data.invalid_json").Key("input").Query(".").HasValue(func() {}).ErrorContains(t, "invalid operation")
}

func TestQueryNotEqual(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("terraform_data.invalid_json").Key("input").Query(".").HasValue(123).AsError(),
		"query result <nil>, for key input not equal to assertion 123",
	)
}

func TestQueryNotExists(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(queryTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	assert.ErrorContains(
		t,
		InPlan(tftest.Plan).That("terraform_data.invalid_json").Key("not_exists").Query(".").HasValue(123).AsError(),
		"key not_exists not found in resource",
	)
}
