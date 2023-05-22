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
	defer test.Cleanup()
	require.NoError(t, err)
	test.Apply().ErrorIsNil(t)
	test.Destroy().ErrorIsNil(t)
}

func TestApplyIdempotent(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	test.ApplyIdempotent().ErrorIsNil(t)
	test.Destroy().ErrorIsNil(t)
}

func TestApplyIdempotentRetryFail(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/applyidempotentretryfail", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	tb := time.Now()
	err = test.ApplyIdempotentRetry(rty).AsError()
	assert.Truef(t, time.Since(tb) >= 10*time.Second, "retry should have waited at least 10 second")
	assert.ErrorContains(t, err, "'terraform plan' unsuccessful after 2 retries")
}

func TestApplyFail(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/applyfail", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	err = test.Apply().AsError()
	assert.ErrorContains(t, err, "test error")
}

func TestApplyIdempotentApplyFail(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/applyfail", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	err = test.ApplyIdempotent().AsError()
	assert.ErrorContains(t, err, "test error")
}

func TestApplyIdempotentRetryApplyFail(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/applyfail", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	err = test.ApplyIdempotentRetry(rty).AsError()
	assert.ErrorContains(t, err, "test error")
}
