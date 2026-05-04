package cmd_test

import (
	"os"
	"testing"

	"github.com/senzing-garage/load/cmd"
	"github.com/stretchr/testify/require"
)

const (
	commandName = "command-name"
	helpFlag    = "--help"
)

// ----------------------------------------------------------------------------
// Test public functions
// ----------------------------------------------------------------------------

func Test_Execute(_ *testing.T) { //nolint
	os.Args = []string{commandName, helpFlag}

	cmd.Execute()
}

func Test_Execute_completion(_ *testing.T) { //nolint
	os.Args = []string{commandName, "completion"}

	cmd.Execute()
}

func Test_Execute_docs(_ *testing.T) { //nolint
	os.Args = []string{commandName, "docs"}

	cmd.Execute()
}

func Test_Execute_help(_ *testing.T) { //nolint
	os.Args = []string{commandName, helpFlag}

	cmd.Execute()
}

func Test_PreRun(_ *testing.T) { //nolint
	args := []string{commandName, helpFlag}
	cmd.PreRun(cmd.RootCmd, args)
}

// func Test_RunE(test *testing.T) {
// 	test.Setenv(
// 		"SENZING_TOOLS_INPUT_URL",
//      "https://raw.githubusercontent.com/senzing-garage/knowledge-base/refs/heads/main/testdata/jsonl/single-customer.jsonl",
// 	)

// 	err := cmd.RunE(cmd.RootCmd, []string{})
// 	require.NoError(test, err)
// }

// func Test_RootCmd(test *testing.T) { //nolint
// 	test.Setenv(
// 		"SENZING_TOOLS_INPUT_URL",
//         "https://raw.githubusercontent.com/senzing-garage/knowledge-base/refs/heads/main/testdata/jsonl/single-customer.jsonl",
// 	)
// 	err := cmd.RootCmd.Execute()
// 	require.NoError(test, err)
// 	err = cmd.RootCmd.RunE(
// 		cmd.RootCmd,
// 		[]string{},
// 	)
// 	require.NoError(test, err)
// }

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
