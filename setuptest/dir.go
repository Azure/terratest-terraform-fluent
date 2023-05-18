package setuptest

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// PrepFunc is a function that is used to prepare the test directory before running Terraform.
// It takes a SetupTestResponse struct as a parameter and returns an error.
//
// Example usage is to create a files in the test directory when testing submodules, e.g. AzureRM provider blocks.
type PrepFunc func(Response) error

// Response is a struct which contains:
//
// - the temporary directory used by the test
// - the plan struct
// - the Terraform options
// - a cleanup function which deletes the temporary directory and provides sanitized serial logging (useful when running parallel tests)
//
// It is returned by the Init* funcs and can be used by the check package.
type Response struct {
	TmpDir  string
	Plan    *terraform.PlanStruct
	Options *terraform.Options
	Cleanup func()
	t       *testing.T
}

// Dirs func begins the fluent test setup process.
// It takes a root directory and a test directory as parameters.
//
// The root directory is the directory containing the terraform code to be tested.
// The test directory is the directory containing the test code,
// it should either be blank to test the code in the root,
// or a relative path beneath the root directory.
//
// Before a Terraform command is run, the code in the root directory will be copied to a temporary directory.
func Dirs(rootdir, testdir string) DirType {
	return DirType{
		RootDir: rootdir,
		TestDir: testdir,
	}
}

// DirType is a type which can be used for more fluent setup of a test.
// It contains the two directories required. The root, and the subdirectory of the root containing the test.
// If the test directory is blank, the test will be run in the root directory.
type DirType struct {
	RootDir string
	TestDir string
}

// WithVars is an method of DirType and allows you to add variables in the form of `map[string]any`.
// It returns a type which can be used for more fluent setup of a test with variables.
func (d DirType) WithVars(vars map[string]any) DirTypeWithVars {
	return DirTypeWithVars{
		RootDir: d.RootDir,
		TestDir: d.TestDir,
		Vars:    vars,
	}
}

// WithVarFiles is an method of DirType and allows you to add variable files in the form of []string.
// It returns a type which can be used for more fluent setup of a test with variable files.
func (d DirType) WithVarFiles(varfiles []string) DirTypeWithVarFiles {
	return DirTypeWithVarFiles{
		RootDir:  d.RootDir,
		TestDir:  d.TestDir,
		VarFiles: varfiles,
	}
}

// DirTypeWithVars is a type which can be used for more fluent setup of a test with variables.
// It is used by the Init* methods to setup the test.
type DirTypeWithVars struct {
	RootDir string
	TestDir string
	Vars    map[string]any
}

// DirTypeWithVarFiles is a type which can be used for more fluent setup of a test with variable files.
// It is used by the Init* methods to setup the test.
type DirTypeWithVarFiles struct {
	RootDir  string
	TestDir  string
	VarFiles []string
}
