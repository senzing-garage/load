package rabbitmq

import (
	"context"
	"fmt"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/go-queueing/queues/rabbitmq"
	"github.com/senzing-garage/go-sdk-abstract-factory/szfactorycreator"
	"github.com/senzing-garage/sz-sdk-go/senzing"
)

// ----------------------------------------------------------------------------

// Read and process records from the given queue until a system interrupt.
func Read(ctx context.Context, urlString, engineConfigJSON, logLevel string, jsonOutput bool) {
	logger = getLogger()

	err := setLogLevel(ctx, logLevel)
	if err != nil {
		panic("Cannot set log level")
	}

	szAbstractFactory, err := szfactorycreator.CreateCoreAbstractFactory(
		"load",
		engineConfigJSON,
		senzing.SzNoLogging,
		senzing.SzInitializeWithDefaultConfiguration,
	)
	if err != nil {
		panic(err)
	}

	defer func() {
		err := szAbstractFactory.Destroy(ctx)
		if err != nil {
			panic(err)
		}
	}()

	szEngine, err := szAbstractFactory.CreateEngine(ctx)
	if err != nil {
		panic(err)
	}

	startErr := rabbitmq.StartManagedConsumer(ctx, urlString, 0, &szEngine, false, logLevel, jsonOutput)

	if startErr != nil {
		log(5000, startErr.Error())
	}

	log(2999)
}

// ----------------------------------------------------------------------------

var (
	logger     logging.Logging
	jsonOutput bool
)

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func getLogger() logging.Logging {
	var err error

	if logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: OptionCallerSkip},
		}

		logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
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
		fmt.Println(fmt.Sprintf(IDMessages[messageNumber], details...)) //nolint
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

	var err error

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return wraperror.Errorf(errForPackage, "invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = getLogger().SetLogLevel(logLevelName)

	return wraperror.Errorf(err, wraperror.NoMessage)
}
