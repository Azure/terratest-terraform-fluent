package setuptest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestDirs(t *testing.T) {
	t.Parallel()

	_, err := Dirs("testdata/depth1", "").WithVars(map[string]interface{}{}).InitPlanShow(t)
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
