package setuptest

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestInit(t *testing.T) {
	t.Parallel()

	t.Run("Vars", func(t *testing.T) {
		t.Parallel()
		v := map[string]any{"test": "testing"}
		test, err := Dirs("testdata/with-vars", "").WithVars(v).Init(t)
		defer test.Cleanup()
		assert.NoError(t, err)
	})

	t.Run("VarFiles", func(t *testing.T) {
		t.Parallel()
		v := []string{"vars.tfvars"}
		test, err := Dirs("testdata/with-vars", "").WithVarFiles(v).Init(t)
		defer test.Cleanup()
		assert.NoError(t, err)
	})

	t.Run("FailVarFiles", func(t *testing.T) {
		t.Parallel()
		v := []string{"vars.tfvars"}
		_, err := Dirs("testdata/notexist", "").WithVarFiles(v).Init(t)
		assert.ErrorContains(t, err, "no such file or directory")
	})

	t.Run("FailVars", func(t *testing.T) {
		t.Parallel()
		v := map[string]any{"test": "testing"}
		_, err := Dirs("testdata/notexist", "").WithVars(v).Init(t)
		assert.ErrorContains(t, err, "no such file or directory")
	})
}
