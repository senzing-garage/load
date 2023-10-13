package load

import (
	"context"
	"fmt"
	"runtime"
	"runtime/debug"
	"sync"
	"time"

	"github.com/senzing/go-logging/logging"
	"github.com/senzing/load/input"
)

// ----------------------------------------------------------------------------
// Types
// ----------------------------------------------------------------------------

type LoadImpl struct {
	EngineConfigJson          string
	EngineLogLevel            int64
	InputURL                  string
	JSONOutput                bool
	logger                    logging.LoggingInterface
	LogLevel                  string
	NumberOfWorkers           int
	MonitoringPeriodInSeconds int
	// RecordMax                 int
	// RecordMin                 int
	RecordMonitor             int
	VisibilityPeriodInSeconds int
}

// ----------------------------------------------------------------------------

// Check at compile time that the implementation adheres to the interface.
var _ Load = (*LoadImpl)(nil)

// ----------------------------------------------------------------------------

func (l *LoadImpl) Load(ctx context.Context) error {

	l.logBuildInfo()
	l.logStats()

	ticker := time.NewTicker(time.Duration(l.MonitoringPeriodInSeconds) * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				l.logStats()
			}
		}
	}()

	return input.Read(ctx, l.InputURL, l.EngineConfigJson, l.EngineLogLevel, l.NumberOfWorkers, l.VisibilityPeriodInSeconds, l.LogLevel, l.JSONOutput)
}

// ----------------------------------------------------------------------------
// Logging --------------------------------------------------------------------
// ----------------------------------------------------------------------------

// Get the Logger singleton.
func (v *LoadImpl) getLogger() logging.LoggingInterface {
	var err error = nil
	if v.logger == nil {
		options := []interface{}{
			&logging.OptionCallerSkip{Value: 4},
		}
		v.logger, err = logging.NewSenzingToolsLogger(ComponentID, IDMessages, options...)
		if err != nil {
			panic(err)
		}
	}
	return v.logger
}

// Log message.
func (v *LoadImpl) log(messageNumber int, details ...interface{}) {
	if v.JSONOutput {
		v.getLogger().Log(messageNumber, details...)
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
func (v *LoadImpl) SetLogLevel(ctx context.Context, logLevelName string) error {
	var err error = nil

	// Verify value of logLevelName.

	if !logging.IsValidLogLevelName(logLevelName) {
		return fmt.Errorf("invalid error level: %s", logLevelName)
	}

	// Set ValidateImpl log level.

	err = v.getLogger().SetLogLevel(logLevelName)
	return err
}

// ----------------------------------------------------------------------------

func (m *LoadImpl) logBuildInfo() {
	buildInfo, ok := debug.ReadBuildInfo()
	if ok {
		m.log(2002, buildInfo.GoVersion, buildInfo.Path, buildInfo.Main.Path, buildInfo.Main.Version)
	} else {
		m.log(3011)
	}
}

// ----------------------------------------------------------------------------

var lock sync.Mutex

func (m *LoadImpl) logStats() {
	lock.Lock()
	defer lock.Unlock()
	cpus := runtime.NumCPU()
	goRoutines := runtime.NumGoroutine()
	cgoCalls := runtime.NumCgoCall()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)
	m.log(2003, cpus, goRoutines, cgoCalls, memStats.NumGC, gcStats.PauseTotal, gcStats.LastGC, memStats.TotalAlloc, memStats.HeapAlloc, memStats.NextGC, memStats.GCSys, memStats.HeapSys, memStats.StackSys, memStats.Sys, memStats.GCCPUFraction)

}
