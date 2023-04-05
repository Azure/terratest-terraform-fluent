package check

import (
	"testing"

	"github.com/Azure/terratest-terraform-fluent/setuptest"
	"github.com/stretchr/testify/require"
)

func TestNumberOfResourcesInPlan(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).NumberOfResourcesEquals(4).ErrorIsNil(t)
}

func TestNumberOfResourcesInPlanWithError(t *testing.T) {
	t.Parallel()

	tftest, err := setuptest.Dirs(basicTestData, "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer tftest.Cleanup()
	InPlan(tftest.Plan).NumberOfResourcesEquals(1).ErrorContains(t, "expected 1 resources, got 4")
}
