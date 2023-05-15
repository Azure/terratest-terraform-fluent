package check

import (
	"testing"

	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/gruntwork-io/terratest/modules/terraform"
	tfjson "github.com/hashicorp/terraform-json"
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

	twkq := mockThatTypeWithKeyQuery(nil, ".")
	err := twkq.HasValue(nil).AsError()
	assert.NoError(t, err)
}

func TestQueryInvalidArgs(t *testing.T) {
	t.Parallel()

	twkq := mockThatTypeWithKeyQuery(nil, ".")

	err := twkq.HasValue(func() {}).AsError()
	assert.ErrorContains(t, err, "invalid operation")
}

func TestQueryNotEqual(t *testing.T) {
	t.Parallel()

	twkq := mockThatTypeWithKeyQuery("{}", "test_key")
	err := twkq.HasValue(123).AsError()
	assert.ErrorContains(t, err, "query result <nil>, for key test_key not equal to assertion 123")
}

func TestQueryNotExists(t *testing.T) {
	t.Parallel()

	twkq := mockThatTypeWithKeyQuery(nil, "not_exists")
	twkq.Key = "not_exists"
	err := twkq.HasValue(nil).AsError()
	assert.ErrorContains(t, err, "key not_exists not found in resource")
}

func TestQueryEscapedJson(t *testing.T) {
	t.Parallel()

	// Create a mock ThatTypeWithKeyQuery object
	val := "{\"properties\":{\"testProperty\":{\"testArray\":[\"testArrayMember\"]},\"testProperty2\":{\"TestArray2\":[]}}}"
	twkq := mockThatTypeWithKeyQuery(val, "properties.testProperty.testArray.0")
	err := twkq.HasValue("testArrayMember")
	assert.NoError(t, err.AsError())
}

func mockThatTypeWithKeyQuery(val any, query string) ThatTypeWithKeyQuery {
	return ThatTypeWithKeyQuery{
		Plan: &terraform.PlanStruct{
			ResourcePlannedValuesMap: map[string]*tfjson.StateResource{
				"test_resource": {
					AttributeValues: map[string]interface{}{
						"test_key": val,
					},
				},
			},
		},
		ResourceName: "test_resource",
		Key:          "test_key",
		Query:        query,
	}
}
