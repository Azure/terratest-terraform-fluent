package setuptest

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/gruntwork-io/terratest/modules/logger"
	"github.com/gruntwork-io/terratest/modules/terraform"
	test_structure "github.com/gruntwork-io/terratest/modules/test-structure"
)

type testExecutor interface {
	Logger() logger.TestLogger
}

var _ testExecutor = executor{}

type executor struct{}

func (executor) Logger() logger.TestLogger {
	l := NewMemoryLogger()
	return l
}

// getDefaultTerraformOptions returns the default Terraform options for the
// given directory.
func getDefaultTerraformOptions(t *testing.T, dir string) *terraform.Options {
	if !strings.HasSuffix(dir, "/") {
		dir += "/"
	}

	o := terraform.Options{
		Logger:       logger.Default,
		PlanFilePath: "tfplan",
		TerraformDir: dir,
		Lock:         true,
		NoColor:      true,
		Vars:         make(map[string]any),
	}
	return terraform.WithDefaultRetryableErrors(t, &o)
}

// CopyTerraformFolderToTempAndCleanUp sets up a temporary copy of the supplied module folder
// It will return a three values which contains the temporary directory, a cleanup function and an error.
//
// The testdir input is the relative path to the test directory, it can be blank if testing the module directly with variables
// or it can be a relative path to the module directory if testing the module using a subdirectory.
//
// Note: This function will only work if the test directory is in a child subdirectory of the test directory.
// e.g. you cannot use parent paths of the moduleDir.
//
// The depth input is used to determine how many directories to go up to make sure we
// fully clean up.
func CopyTerraformFolderToTempAndCleanUp(t *testing.T, moduleDir string, testDir string) (string, func() error, error) {
	//var resp Response
	//resp.t = t
	tmp := test_structure.CopyTerraformFolderToTemp(t, moduleDir, testDir)
	// We normalise, then work out the depth of the test directory relative
	// to the test so we know how many/ directories to go up to get to the root.
	// We can then delete the right directory when cleaning up.
	//resp.TmpDir = tmp

	absTestPath := filepath.Join(moduleDir, testDir)
	relPath, err := filepath.Rel(moduleDir, absTestPath)
	if err != nil {
		err = fmt.Errorf("could not get relative path to test directory: %v", err)
		return "", nil, err
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

	f := func() error {
		err := os.RemoveAll(dir)
		return err
	}

	return tmp, f, nil
}

// setup performs the copying of the module dirs to a tmp location
// and returns a Response struct and an error.
func setup(t *testing.T, moduleDir, testDir string, prep PrepFunc) (Response, error) {
	resp := Response{}
	subdir := filepath.Join(moduleDir, testDir)
	_, err := os.Stat(subdir)
	if os.IsNotExist(err) {
		return resp, err
	}
	resp.t = t
	tmp, cleanup, err := CopyTerraformFolderToTempAndCleanUp(t, moduleDir, testDir)
	if err != nil {
		return resp, err
	}
	resp.TmpDir = tmp
	resp.Options = getDefaultTerraformOptions(t, tmp)

	if prep != nil {
		err = prep(resp)
		if err != nil {
			return resp, err
		}
	}

	l := testExecutor(executor{}).Logger()
	resp.Options.Logger = logger.New(l)
	funcs := []func() error{cleanup}
	c, ok := l.(io.Closer)
	if ok {
		funcs = append(funcs, c.Close)
	}
	resp.Cleanup = func() {
		for _, fn := range funcs {
			_ = fn()
		}
	}
	return resp, nil
}
