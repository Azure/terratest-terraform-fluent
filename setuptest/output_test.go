package setuptest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutput(t *testing.T) {
	v := make(map[string]any)
	v["test"] = "test"
	test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
	defer test.Cleanup()
	require.NoError(t, err)
	err = test.Apply().AsError()
	assert.NoError(t, err)
	err = test.Output("test").HasValue("test").AsError()
	assert.NoError(t, err)
}
