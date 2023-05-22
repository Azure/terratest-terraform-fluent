package setuptest

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestOutput(t *testing.T) {
	t.Parallel()

	t.Run("String", func(t *testing.T) {
		t.Parallel()
		v := make(map[string]any)
		v["test_string"] = "test"
		test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
		defer test.Cleanup()
		require.NoError(t, err)
		err = test.Apply().AsError()
		assert.NoError(t, err)
		err = test.Output("test_string").HasValue("test").AsError()
		assert.NoError(t, err)
	})

	t.Run("Number", func(t *testing.T) {
		t.Parallel()
		v := make(map[string]any)
		v["test_number"] = 2
		test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
		defer test.Cleanup()
		require.NoError(t, err)
		err = test.Apply().AsError()
		assert.NoError(t, err)
		err = test.Output("test_number").HasValue(2).AsError()
		assert.NoError(t, err)
	})

	t.Run("Map", func(t *testing.T) {
		t.Parallel()
		v := make(map[string]any)
		v["test_map"] = map[string]any{
			"test_key": "test_value",
		}
		test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
		defer test.Cleanup()
		require.NoError(t, err)
		err = test.Apply().AsError()
		assert.NoError(t, err)
		err = test.Output("test_map").Query("test_key").HasValue("test_value").AsError()
		assert.NoError(t, err)
		err = test.Output("test_map").HasValue(v["test_map"]).AsError()
		assert.NoError(t, err)
	})

	t.Run("Bool", func(t *testing.T) {
		t.Parallel()
		v := make(map[string]any)
		v["test_bool"] = true
		test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
		defer test.Cleanup()
		require.NoError(t, err)
		err = test.Apply().AsError()
		assert.NoError(t, err)
		err = test.Output("test_bool").HasValue(true).AsError()
		assert.NoError(t, err)
	})

	t.Run("List", func(t *testing.T) {
		t.Parallel()
		v := make(map[string]any)
		v["test_list"] = []any{"test_value", "test_value2"}
		test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
		defer test.Cleanup()
		require.NoError(t, err)
		err = test.Apply().AsError()
		assert.NoError(t, err)
		err = test.Output("test_list").HasValue(v["test_list"]).AsError()
		assert.NoError(t, err)
		err = test.Output("test_list").Query("#").HasValue(2).AsError()
		assert.NoError(t, err)
	})

	t.Run("Set", func(t *testing.T) {
		t.Parallel()
		v := make(map[string]any)
		v["test_set"] = []any{"test_value", "test_value2"}
		test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
		defer test.Cleanup()
		require.NoError(t, err)
		err = test.Apply().AsError()
		assert.NoError(t, err)
		err = test.Output("test_set").HasValue(v["test_set"]).AsError()
		assert.NoError(t, err)
		err = test.Output("test_set").Query("#").HasValue(2).AsError()
		assert.NoError(t, err)
	})
}

func TestOutputNotFound(t *testing.T) {
	t.Parallel()
	v := make(map[string]any)
	test, err := Dirs("testdata/output", "").WithVars(v).Init(t)
	defer test.Cleanup()
	require.NoError(t, err)
	err = test.Apply().AsError()
	assert.NoError(t, err)
	o := test.Output("not_found")
	assert.False(t, o.Exist)
}
