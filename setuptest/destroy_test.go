package setuptest

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDestroyRetry(t *testing.T) {
	t.Parallel()

	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	test.ApplyIdempotent().ErrorIsNil(t)
	test.DestroyRetry(rty).ErrorIsNil(t)
}

func TestDestroy(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/depth1", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	test.ApplyIdempotent().ErrorIsNil(t)
	test.Destroy().ErrorIsNil(t)
}

func TestDestroyFail(t *testing.T) {
	t.Parallel()

	test, err := Dirs("testdata/destroyfailretryok", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	test.Apply().ErrorIsNil(t)
	err = test.Destroy().AsError()
	assert.ErrorContains(t, err, "error while running command: exit status 1")
}

func TestDestroyRetryOnceFail(t *testing.T) {
	rty := Retry{
		Max:  1,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/destroyfailretryok", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	test.Apply().ErrorIsNil(t)
	err = test.DestroyRetry(rty).AsError()
	assert.ErrorContains(t, err, "error while running command: exit status 1")
}

func TestDestroyRetryOnce(t *testing.T) {
	rty := Retry{
		Max:  2,
		Wait: time.Second * 10,
	}
	test, err := Dirs("testdata/destroyfailretryok", "").WithVars(nil).InitPlanShow(t)
	defer test.Cleanup()
	require.NoError(t, err)
	test.Apply().ErrorIsNil(t)
	tb := time.Now()
	test.DestroyRetry(rty).ErrorIsNil(t)
	assert.Truef(t, time.Since(tb) >= 10*time.Second, "retry should have waited at least 10 second")
}
