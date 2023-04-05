package setuptest

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

// getDefaultTerraformOptions returns the default Terraform options for the
// given directory.
func getDefaultTerraformOptions(t *testing.T, dir string) *terraform.Options {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}
	o := terraform.Options{
		Logger:       logger.TestingT,
		PlanFilePath: "tfplan",
		TerraformDir: dir,
		Lock:         true,
		NoColor:      true,

		Vars: make(map[string]interface{}),
	}
	return terraform.WithDefaultRetryableErrors(t, &o)
}

// CopyTerraformFolderToTempAndCleanUp sets up a temporary copy of the supplied module folder
// It will return a SetupTestResponse struct which contains the temporary directory, a cleanup function and an error.
//
// The testdir input is the relative path to the test directory, it can be blank if testing the module directly with variables
// or it can be a relative path to the module directory if testing the module using a subdirectory.
//
// Note: This function will only work if the test directory is in a child subdirectory of the test directory.
// e.g. you cannot use parent paths of the moduleDir.
//
// The depth input is used to determine how many directories to go up to make sure we
// fully clean up.
//
// The function will return the temporary directory to use with the terraform options struct, as well as
// a function that can be used with defer to clean up afterwards.
func CopyTerraformFolderToTempAndCleanUp(t *testing.T, moduleDir string, testDir string) (Response, error) {
	var resp Response
	tmp := test_structure.CopyTerraformFolderToTemp(t, moduleDir, testDir)
	// We normalise, then work out the depth of the test directory relative
	// to the test so we know how many/ directories to go up to get to the root.
	// We can then delete the right directory when cleaning up.
	resp.TmpDir = tmp

	absTestPath := filepath.Join(moduleDir, testDir)
	relPath, err := filepath.Rel(moduleDir, absTestPath)
	if err != nil {
		err = fmt.Errorf("could not get relative path to test directory: %v", err)
		return resp, err
	}
	list := strings.Split(relPath, string(os.PathSeparator))
	depth := len(list)
	if len(list) > 1 || list[0] != "." {
		depth++
	}
	dir := tmp
	for i := 0; i < depth; i++ {
		dir = filepath.Dir(dir)
	}

	resp.Cleanup = func() {
		removeTestDir(t, dir)
	}

	return resp, nil
}

// removeTestDir removes the supplied test directory
func removeTestDir(t *testing.T, dir string) {
	err := os.RemoveAll(dir)
	if err != nil {
		t.Logf("error removing test directory %s: %v", dir, err)
	}
}
