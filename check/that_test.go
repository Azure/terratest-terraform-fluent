package check

import (
	"testing"

	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	basicTestData = "testdata/basic"
)

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
