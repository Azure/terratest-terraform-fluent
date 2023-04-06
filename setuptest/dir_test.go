package setuptest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestDirs(t *testing.T) {
	t.Parallel()

	_, err := Dirs("testdata/depth1", "").WithVars(map[string]interface{}{}).InitPlanShow(t)
	require.NoError(t, err)
}

func TestDirsWithVars(t *testing.T) {
	t.Parallel()
	vars := map[string]interface{}{
		"test": "testing",
	}
	tftest, err := Dirs("testdata/with-vars", "").WithVars(vars).InitPlanShow(t)
	defer tftest.Cleanup()
	assert.Equal(t, "testing", tftest.Plan.RawPlan.OutputChanges["test"].After)
	require.NoError(t, err)
}

func TestDirsWithVarFiles(t *testing.T) {
	t.Parallel()
	vf := []string{"vars.tfvars"}
	_, err := Dirs("testdata/with-vars", "").WithVarFiles(vf).InitPlanShow(t)
	require.NoError(t, err)
}

func TestDirsWithVarFilesWithFunc(t *testing.T) {
	t.Parallel()

	var f PrepFunc = func(resp Response) error {
		f, err := os.Create(filepath.Join(resp.TmpDir, "test.txt"))
		if err != nil {
			return err
		}
		return f.Close()
	}

	vf := []string{"vars.tfvars"}
	tftest, err := Dirs("testdata/with-vars", "").WithVarFiles(vf).InitPlanShowWithPrepFunc(t, f)
	defer tftest.Cleanup()
	assert.Equal(t, "testing", tftest.Plan.RawPlan.OutputChanges["test"].After)
	require.NoError(t, err)
}

func TestDirsWithFunc(t *testing.T) {
	t.Parallel()

	var f PrepFunc = func(resp Response) error {
		f, err := os.Create(filepath.Join(resp.TmpDir, "test.txt"))
		if err != nil {
			return err
		}
		return f.Close()
	}

	test, err := Dirs("testdata/depth1", "").WithVars(map[string]interface{}{}).InitPlanShowWithPrepFunc(t, f)
	require.NoError(t, err)
	require.FileExists(t, filepath.Join(test.TmpDir, "test.txt"))
}

func TestDirsNotExist(t *testing.T) {
	t.Parallel()

	_, err := Dirs("testdata/notexist", "").WithVars(nil).InitPlanShow(t)

	require.True(t, os.IsNotExist(err))
}

func TestDirsWithFuncNotExist(t *testing.T) {
	t.Parallel()

	_, err := Dirs("testdata/notexist", "").WithVars(nil).InitPlanShowWithPrepFunc(t, nil)

	require.True(t, os.IsNotExist(err))
}

func TestDirsWithVarFilesNotExist(t *testing.T) {
	t.Parallel()

	_, err := Dirs("testdata/notexist", "").WithVarFiles(nil).InitPlanShow(t)

	require.True(t, os.IsNotExist(err))
}

func TestDirsWithVarFilesWithFuncNotExist(t *testing.T) {
	t.Parallel()

	_, err := Dirs("testdata/notexist", "").WithVarFiles(nil).InitPlanShowWithPrepFunc(t, nil)

	require.True(t, os.IsNotExist(err))
}
