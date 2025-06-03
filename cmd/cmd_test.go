package cmd_test

import (
	"os"
	"testing"

	"github.com/senzing-garage/load/cmd"
	"github.com/stretchr/testify/require"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name", "--help"}

	cmd.Execute()
}

func Test_Execute_completion(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name", "completion"}

	cmd.Execute()
}

func Test_Execute_docs(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name", "docs"}

	cmd.Execute()
}

func Test_Execute_help(test *testing.T) {
	test.Parallel()

	os.Args = []string{"command-name", "--help"}

	cmd.Execute()
}

func Test_PreRun(test *testing.T) {
	test.Parallel()

	args := []string{"command-name", "--help"}
	cmd.PreRun(cmd.RootCmd, args)
}

func Test_CompletionCmd(test *testing.T) {
	test.Parallel()

	err := cmd.CompletionCmd.Execute()
	require.NoError(test, err)
	err = cmd.CompletionCmd.RunE(cmd.CompletionCmd, []string{})
	require.NoError(test, err)
}

func Test_DocsCmd(test *testing.T) {
	test.Parallel()

	err := cmd.DocsCmd.Execute()
	require.NoError(test, err)
	err = cmd.DocsCmd.RunE(cmd.DocsCmd, []string{})
	require.NoError(test, err)
}
