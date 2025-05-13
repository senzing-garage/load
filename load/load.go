package load

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/senzing-garage/go-helpers/wraperror"
	"github.com/senzing-garage/go-logging/logging"
	"github.com/senzing-garage/load/input"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type BasicLoad struct {
	EngineConfigJSON          string
	EngineLogLevel            int64
	InputURL                  string
	JSONOutput                bool
	logger                    logging.Logging
	LogLevel                  string
	NumberOfWorkers           int
	MonitoringPeriodInSeconds int
	// RecordMax                 int
	// RecordMin                 int
	RecordMonitor             int
	VisibilityPeriodInSeconds int32
}

// ----------------------------------------------------------------------------

// Check at compile time that the implementation adheres to the interface.
var _ Load = (*BasicLoad)(nil)

// ----------------------------------------------------------------------------
// Interface methods
// ----------------------------------------------------------------------------

func (load *BasicLoad) Load(ctx context.Context) error {

	load.logBuildInfo()
	load.logStats()

	ticker := time.NewTicker(time.Duration(load.MonitoringPeriodInSeconds) * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				load.logStats()
			}
		}
	}()

	return input.Read(
		ctx,
		load.InputURL,
		load.EngineConfigJSON,
		load.EngineLogLevel,
		load.NumberOfWorkers,
		load.VisibilityPeriodInSeconds,
		load.LogLevel,
		load.JSONOutput,
	)
}

/*
The SetLogLevel method sets the level of logging.

Input
  - ctx: A context to control lifecycle.
  - logLevel: The desired log level. TRACE, DEBUG, INFO, WARN, ERROR, FATAL or PANIC.
*/
func (load *BasicLoad) SetLogLevel(ctx context.Context, logLevelName string) error {
	_ = ctx
	var err error

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return wraperror.Errorf(errForPackage, "invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = load.getLogger().SetLogLevel(logLevelName)
	return err
}

// ----------------------------------------------------------------------------
// Private methods
// ----------------------------------------------------------------------------

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (load *BasicLoad) getLogger() logging.Logging {
	var err error
	if load.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: OptionCallerSkip},
		}
		load.logger, err = logging.NewSenzingLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return load.logger
}

// Log message.
func (load *BasicLoad) log(messageNumber int, details ...interface{}) {
	if load.JSONOutput {
		load.getLogger().Log(messageNumber, details...)
	} else {
		fmt.Println(fmt.Sprintf(IDMessages[messageNumber], details...)) //nolint
	}
}

// ----------------------------------------------------------------------------

func (load *BasicLoad) logBuildInfo() {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		load.log(2002, buildInfo.GoVersion, buildInfo.Path, buildInfo.Main.Path, buildInfo.Main.Version)
	} else {
		load.log(3011)
	}
}

// ----------------------------------------------------------------------------

var lock sync.Mutex

func (load *BasicLoad) logStats() {
	lock.Lock()
	defer lock.Unlock()
	cpus := runtime.NumCPU()
	goRoutines := runtime.NumGoroutine()
	cgoCalls := runtime.NumCgoCall()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)
	load.log(
		2003,
		cpus,
		goRoutines,
		cgoCalls,
		memStats.NumGC,
		gcStats.PauseTotal,
		gcStats.LastGC,
		memStats.TotalAlloc,
		memStats.HeapAlloc,
		memStats.NextGC,
		memStats.GCSys,
		memStats.HeapSys,
		memStats.StackSys,
		memStats.Sys,
		memStats.GCCPUFraction,
	)

}
