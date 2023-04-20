package setuptest

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
)

// PrepFunc is a function that is used to prepare the test directory before running Terraform.
// It takes a SetupTestResponse struct as a parameter and returns an error.
//
// Example usage is to create a files in the test directory when testing submodules, e.g. AzureRM provider blocks.
type PrepFunc func(Response) error

// Response is a struct which contains the temporary directory, the plan, the Terraform options and a cleanup function.
// It is returned by the InitPlanShow funcs and can be used by the check package.
type Response struct {
	TmpDir  string
	Plan    *terraform.PlanStruct
	Options *terraform.Options
	Cleanup func()
}

// Dirs func begins the fluent test testup process.
// It takes a root directory and a test directory as parameters.
//
// The root directory is the directory containing the terraform code to be tested.
//
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

// DirType is a type which can be used for more fluent setup of a test
type DirType struct {
	RootDir string
	TestDir string
}

// WithVars allows you to add variables in the form of `map[string]any`.
// It returns a type which can be used for more fluent setup of a test with variables.
func (d DirType) WithVars(vars map[string]any) DirTypeWithVars {
	return DirTypeWithVars{
		RootDir: d.RootDir,
		TestDir: d.TestDir,
		Vars:    vars,
	}
}

// WithVarFiles returns a type which can be used for more fluent setup of a test with variable files
func (d DirType) WithVarFiles(varfiles []string) DirTypeWithVarFiles {
	return DirTypeWithVarFiles{
		RootDir:  d.RootDir,
		TestDir:  d.TestDir,
		VarFiles: varfiles,
	}
}

// DirTypeWithVars is a type which can be used for more fluent setup of a test with variables
type DirTypeWithVars struct {
	RootDir string
	TestDir string
	Vars    map[string]any
}

// DirTypeWithVarFiles is a type which can be used for more fluent setup of a test with variable files
type DirTypeWithVarFiles struct {
	RootDir  string
	TestDir  string
	VarFiles []string
}

// InitPlanShow is a wrapper around terraform.InitAndPlanAndShowWithStructE
// It takes a test object as a parameter and returns a SetupTestResponse.
//
// The SetupTestResponse contains the temporary directory, the plan, the Terraform options and a cleanup function.
// The temporary directory is the directory containing a copy of the code specified by the Dirs func.
// The plan is the plan struct generated by terraform, which can be used by the check package.
// The cleanup function is a function which should be used with defer to clean up the temporary directory.
func (dtv DirTypeWithVars) InitPlanShow(t *testing.T) (Response, error) {
	subdir := filepath.Join(dtv.RootDir, dtv.TestDir)
	_, err := os.Stat(subdir)
	if os.IsNotExist(err) {
		return Response{}, err
	}
	// copy test to tmp dir
	resp, err := CopyTerraformFolderToTempAndCleanUp(t, dtv.RootDir, dtv.TestDir)
	if err != nil {
		return resp, err
	}

	// Run terraform
	opts := getDefaultTerraformOptions(t, resp.TmpDir)
	opts.Vars = dtv.Vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)

	resp.Options = opts
	resp.Plan = plan
	return resp, err
}

// InitPlanShowWithPrepFunc is a wrapper around terraform.InitAndPlanAndShowWithStructE
// It takes a test object and a SetupTestPrepFunc as a parameter and returns a SetupTestResponse.
// The PrepFunc is executed after the test has been copied to a tmp directory,
// allowing file modifications to be made before running Terraform.
//
// The SetupTestResponse contains the temporary directory, the plan, the Terraform options and a cleanup function.
// The temporary directory is the directory containing a copy of the code specified by the Dirs func.
// The terraform options are the options used to run terraform and can be used by the apply functions.
// The plan is the plan struct generated by terraform, which can be used by the check package.
// The cleanup function is a function which should be used with defer to clean up the temporary directory.
func (dtv DirTypeWithVars) InitPlanShowWithPrepFunc(t *testing.T, f PrepFunc) (Response, error) {
	subdir := filepath.Join(dtv.RootDir, dtv.TestDir)
	_, err := os.Stat(subdir)
	if os.IsNotExist(err) {
		return Response{}, err
	}

	// Copy test to tmp dir
	resp, err := CopyTerraformFolderToTempAndCleanUp(t, dtv.RootDir, dtv.TestDir)
	if err != nil {
		return resp, err
	}

	// Run the prep function
	if err := f(resp); err != nil {
		return resp, err
	}

	// Run terraform
	opts := getDefaultTerraformOptions(t, resp.TmpDir)
	opts.Vars = dtv.Vars
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)

	resp.Options = opts
	resp.Plan = plan
	return resp, err
}

// InitPlanShow is a wrapper around terraform.InitAndPlanAndShowWithStructE
// It takes a test object as a parameter and returns a SetupTestResponse.
//
// The SetupTestResponse contains the temporary directory, the plan, the Terraform options and a cleanup function.
// The temporary directory is the directory containing a copy of the code specified by the Dirs func.
// The plan is the plan struct generated by terraform, which can be used by the check package.
// The cleanup function is a function which should be used with defer to clean up the temporary directory.
func (dtvf DirTypeWithVarFiles) InitPlanShow(t *testing.T) (Response, error) {
	subdir := filepath.Join(dtvf.RootDir, dtvf.TestDir)
	_, err := os.Stat(subdir)
	if os.IsNotExist(err) {
		return Response{}, err
	}
	// copy test to tmp dir
	resp, err := CopyTerraformFolderToTempAndCleanUp(t, dtvf.RootDir, dtvf.TestDir)
	if err != nil {
		return resp, err
	}

	// Run terraform
	opts := getDefaultTerraformOptions(t, resp.TmpDir)
	opts.VarFiles = dtvf.VarFiles
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)

	resp.Options = opts
	resp.Plan = plan
	return resp, err
}

// InitPlanShow is a wrapper around terraform.InitAndPlanAndShowWithStructE
// It takes a test object and a SetupTestPrepFunc as a parameter and returns a SetupTestResponse.
// The PrepFunc is executed after the test has been coped to a tmp directory,
// allowing file modifications to be made before running terraform.
//
// The SetupTestResponse contains the temporary directory, the plan, the Terraform options and a cleanup function.
// The temporary directory is the directory containing a copy of the code specified by the Dirs func.
// The terraform options are the options used to run terraform and can be used by the apply functions.
// The plan is the plan struct generated by terraform, which can be used by the check package.
// The cleanup function is a function which should be used with defer to clean up the temporary directory.
func (dtvf DirTypeWithVarFiles) InitPlanShowWithPrepFunc(t *testing.T, f PrepFunc) (Response, error) {
	subdir := filepath.Join(dtvf.RootDir, dtvf.TestDir)
	_, err := os.Stat(subdir)
	if os.IsNotExist(err) {
		return Response{}, err
	}

	// Copy test to tmp dir
	resp, err := CopyTerraformFolderToTempAndCleanUp(t, dtvf.RootDir, dtvf.TestDir)
	if err != nil {
		return resp, err
	}

	// Run the prep function
	if err := f(resp); err != nil {
		return resp, err
	}

	// Run terraform
	opts := getDefaultTerraformOptions(t, resp.TmpDir)
	opts.VarFiles = dtvf.VarFiles
	plan, err := terraform.InitAndPlanAndShowWithStructE(t, opts)

	resp.Options = opts
	resp.Plan = plan
	return resp, err
}
