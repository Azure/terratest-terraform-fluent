package setuptest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestApply(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.Apply(t).ErrorIsNil(t)
	test.Destroy(t).ErrorIsNil(t)
}

func TestApplyIdempotent(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.ApplyIdempotent(t).ErrorIsNil(t)
	test.Destroy(t).ErrorIsNil(t)
}

func TestApplyIdempotentRetryFail(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/applyidempotentretryfail", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer test.Destroy(t) //nolint:errcheck
	tb := time.Now()
	err = test.ApplyIdempotentRetry(t, rty).AsError()
	assert.Truef(t, time.Since(tb) >= 20*time.Second, "retry should have waited at least 20 second")
	assert.ErrorContains(t, err, "terraform configuration not idempotent")
}

func TestApplyFail(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/applyfail", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer test.Destroy(t) //nolint:errcheck
	err = test.Apply(t).AsError()
	assert.ErrorContains(t, err, "test error")
}

func TestApplyIdempotentApplyFail(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/applyfail", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer test.Destroy(t) //nolint:errcheck
	err = test.ApplyIdempotent(t).AsError()
	assert.ErrorContains(t, err, "test error")
}

func TestApplyIdempotentRetryApplyFail(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/applyfail", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	defer test.Destroy(t) //nolint:errcheck
	err = test.ApplyIdempotentRetry(t, rty).AsError()
	assert.ErrorContains(t, err, "test error")
}

func TestDestroyRetry(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	require.NoError(t, err)
	test.ApplyIdempotent(t).ErrorIsNil(t)
	test.DestroyRetry(t, rty).ErrorIsNil(t)
}
