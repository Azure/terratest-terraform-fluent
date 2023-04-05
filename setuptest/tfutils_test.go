package setuptest

import (
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestCopyFilesToTempAndCleanupDepth1 tests that the parent of the temp directory
// is cleaned up after the cleanup function is called.
//
// The terraform.CopyTerraformFolderToTemp function creates a temporary directory in the os.TempDir
// called t.Name(){RANDOMNUMBERS}/{moduleDir}/{testDir}
// Therefore we need to make sure the parent of {moduledir} is cleaned up.
func TestCopyFilesToTempAndCleanupDepth1(t *testing.T) {
	t.Parallel()

	resp, err := CopyTerraformFolderToTempAndCleanUp(t, "testdata/depth1", "")
	assert.NoError(t, err)
	assert.DirExists(t, resp.TmpDir)
	resp.Cleanup()
	parent := filepath.Dir(resp.TmpDir)
	t.Logf("parent: %s", parent)
	assert.NoDirExists(t, parent)
}

func TestCopyFilesToTempAndCleanupDepth2(t *testing.T) {
	t.Parallel()

	resp, err := CopyTerraformFolderToTempAndCleanUp(t, "testdata/depth2", "subdir")
	assert.NoError(t, err)
	assert.DirExists(t, resp.TmpDir)
	resp.Cleanup()
	parent := filepath.Dir(filepath.Dir(resp.TmpDir))
	t.Logf("parent: %s", parent)
	assert.NoDirExists(t, parent)
}

func TestCopyFilesToTempAndCleanupDepth3(t *testing.T) {
	t.Parallel()

	resp, err := CopyTerraformFolderToTempAndCleanUp(t, "testdata/depth3", "subdir/subdir2")
	assert.NoError(t, err)
	assert.DirExists(t, resp.TmpDir)
	resp.Cleanup()
	parent := filepath.Dir(filepath.Dir(filepath.Dir(resp.TmpDir)))
	t.Logf("parent: %s", parent)
	assert.NoDirExists(t, parent)
}
