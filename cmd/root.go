/*
 */
package cmd

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/senzing/go-cmdhelping/cmdhelper"
	"github.com/senzing/go-cmdhelping/option"
	"github.com/senzing/load/load"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	Short string = "Load records into Senzing."
	Use   string = "load"
	Long  string = `
	Welcome to load!
	This tool will load records into Senzing. It validates the records conform to the Generic Entity Specification.

	For example:

	load --input-url "amqp://guest:guest@192.168.6.96:5672"
	load --input-url "https://public-read-access.s3.amazonaws.com/TestDataSets/SenzingTruthSet/truth-set-3.0.0.jsonl"
    `
)

// ----------------------------------------------------------------------------
// Context variables
// ----------------------------------------------------------------------------

var ContextVariables = []option.ContextVariable{
	option.DelayInSeconds,
	option.EngineConfigurationJson,
	option.EngineModuleName.SetDefault(fmt.Sprintf("load-%d", time.Now().Unix())),
	option.InputFileType,
	option.InputURL,
	option.JSONOutput,
	option.LogLevel,
	option.MonitoringPeriodInSeconds,
	option.OutputURL,
	option.RecordMax,
	option.RecordMin,
	option.RecordMonitor,
}

// ----------------------------------------------------------------------------
// Private functions
// ----------------------------------------------------------------------------

// Since init() is always invoked, define command line parameters.
func init() {
	cmdhelper.Init(RootCmd, ContextVariables)
}

// ----------------------------------------------------------------------------
// Public functions
// ----------------------------------------------------------------------------

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the RootCmd.
func Execute() {
	err := RootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// Used in construction of cobra.Command
func PreRun(cobraCommand *cobra.Command, args []string) {
	cmdhelper.PreRun(cobraCommand, args, Use, ContextVariables)
}

// Used in construction of cobra.Command
func RunE(_ *cobra.Command, _ []string) error {
	jsonOutput := viper.GetBool(option.JSONOutput.Arg)
	if !jsonOutput {
		fmt.Println("Run with the following parameters:")
		for _, key := range viper.AllKeys() {
			fmt.Println("  - ", key, " = ", viper.Get(key))
		}
	}

	if viper.GetInt(option.DelayInSeconds.Arg) > 0 {
		if !jsonOutput {
			fmt.Println(time.Now(), "Sleep for", viper.GetInt(option.DelayInSeconds.Arg), "seconds to let queues and databases settle down and come up.")
		}
		time.Sleep(time.Duration(viper.GetInt(option.DelayInSeconds.Arg)) * time.Second)
	}

	ctx := context.Background()

	loader := &load.LoadImpl{
		EngineConfigJson: viper.GetString(option.EngineConfigurationJson.Arg),
		// FileType:                  viper.GetString(option.InputFileType.Arg),
		InputURL:                  viper.GetString(option.InputURL.Arg),
		JSONOutput:                viper.GetBool(option.JSONOutput.Arg),
		LogLevel:                  viper.GetString(option.LogLevel.Arg),
		MonitoringPeriodInSeconds: viper.GetInt(option.MonitoringPeriodInSeconds.Arg),
		// RecordMax:                 viper.GetInt(option.RecordMax),
		// RecordMin:                 viper.GetInt(option.RecordMin),
		RecordMonitor: viper.GetInt(option.RecordMonitor.Arg),
	}
	return loader.Load(ctx)
}

// Used in construction of cobra.Command
func Version() string {
	return cmdhelper.Version(githubVersion, githubIteration)
}

// ----------------------------------------------------------------------------
// Command
// ----------------------------------------------------------------------------

// RootCmd represents the command.
var RootCmd = &cobra.Command{
	Use:     Use,
	Short:   Short,
	Long:    Long,
	PreRun:  PreRun,
	RunE:    RunE,
	Version: Version(),
}
