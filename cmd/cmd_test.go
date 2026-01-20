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

func Test_CompletionCmd(test *testing.T) { //nolint
	err := cmd.CompletionCmd.Execute()
	require.NoError(test, err)
	err = cmd.CompletionCmd.RunE(cmd.CompletionCmd, []string{})
	require.NoError(test, err)
}

func Test_DocsCmd(test *testing.T) { //nolint
	err := cmd.DocsCmd.Execute()
	require.NoError(test, err)
	err = cmd.DocsCmd.RunE(cmd.DocsCmd, []string{})
	require.NoError(test, err)
}

func Test_Execute(_ *testing.T) { //nolint
	os.Args = []string{"command-name", "--help"}

	cmd.Execute()
}

func Test_Execute_completion(_ *testing.T) { //nolint
	os.Args = []string{"command-name", "completion"}

	cmd.Execute()
}

func Test_Execute_docs(_ *testing.T) { //nolint
	os.Args = []string{"command-name", "docs"}

	cmd.Execute()
}

func Test_Execute_help(_ *testing.T) { //nolint
	os.Args = []string{"command-name", "--help"}

	cmd.Execute()
}

func Test_PreRun(_ *testing.T) { //nolint
	args := []string{"command-name", "--help"}
	cmd.PreRun(cmd.RootCmd, args)
}

func Test_RunE(test *testing.T) {
	test.Setenv("SENZING_TOOLS_AVOID_SERVING", "true")

	err := cmd.RunE(cmd.RootCmd, []string{})
	require.NoError(test, err)
}

func Test_RootCmd(test *testing.T) { //nolint
	err := cmd.RootCmd.Execute()
	require.NoError(test, err)
	err = cmd.RootCmd.RunE(cmd.RootCmd, []string{})
	require.NoError(test, err)
}
