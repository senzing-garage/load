package sqs

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-queueing/queues/sqs"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	"github.com/senzing-garage/sz-sdk-go/sz"
)

// ----------------------------------------------------------------------------

// read and process records from the given queue until a system interrupt
func Read(ctx context.Context, urlString, engineConfigJson string, engineLogLevel int64, numberOfWorkers, visibilityPeriodInSeconds int, logLevel string, jsonOutput bool) {

	logger = getLogger()
	err := setLogLevel(ctx, logLevel)
	if err != nil {
		panic("Cannot set log level")
	}

	// Work with szEngine.
	szEngine := createG2Engine(ctx, engineConfigJson, engineLogLevel)
	defer szEngine.Destroy(ctx)

	startErr := sqs.StartManagedConsumer(ctx, urlString, numberOfWorkers, szEngine, false, int32(visibilityPeriodInSeconds), logLevel, jsonOutput)
	if startErr != nil {
		log(5000, startErr.Error())
	}
	log(2999)
}

// ----------------------------------------------------------------------------

func getAbstractFactory(ctx context.Context, engineConfigJson string, verboseLogging int64) sz.SzAbstractFactory {
	_ = ctx
	result, err := szfactorycreator.CreateCoreAbstractFactory("load", engineConfigJson, verboseLogging, sz.SZ_INITIALIZE_WITH_DEFAULT_CONFIGURATION)
	if err != nil {
		panic(err)
	}
	return result
}

// create a G2Engine object, on error this function panics.
// see failOnError
func createG2Engine(ctx context.Context, settings string, verboseLogging int64) sz.SzEngine {
	result, err := getAbstractFactory(ctx, settings, verboseLogging).CreateSzEngine(ctx)
	if err != nil {
		log(2004, err.Error())
	}
	return result
}

var logger logging.LoggingInterface = nil
var jsonOutput bool = false

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func getLogger() logging.LoggingInterface {
	var err error = nil
	if logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		logger, err = logging.NewSenzingToolsLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return logger
}

// Log message.
func log(messageNumber int, details ...interface{}) {
	if jsonOutput {
		getLogger().Log(messageNumber, details...)
	} else {
		fmt.Println(fmt.Sprintf(IDMessages[messageNumber], details...))
	}
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func setLogLevel(ctx context.Context, logLevelName string) error {
	_ = ctx
	var err error = nil

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = getLogger().SetLogLevel(logLevelName)
	return err
}
