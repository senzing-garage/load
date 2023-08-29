package load

import (
	"context"
	"fmt"
	"os/exec"
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
	EngineLogLevel            int
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

	logOSInfo()
	logBuildInfo()
	logStats()

	ticker := time.NewTicker(time.Duration(l.MonitoringPeriodInSeconds) * time.Second)
	go func() {
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				logStats()
			}
		}
	}()

	return input.Read(ctx, l.InputURL, l.EngineConfigJson, l.EngineLogLevel, l.NumberOfWorkers, l.VisibilityPeriodInSeconds, l.LogLevel, l.JSONOutput)
}

// ----------------------------------------------------------------------------

func logBuildInfo() {
	buildInfo, ok := debug.ReadBuildInfo()
	fmt.Println("---------------------------------------------------------------")
	if ok {
		fmt.Println("GoVersion:", buildInfo.GoVersion)
		fmt.Println("Path:", buildInfo.Path)
		fmt.Println("Main.Path:", buildInfo.Main.Path)
		fmt.Println("Main.Version:", buildInfo.Main.Version)
	} else {
		fmt.Println("Unable to read build info.")
	}
}

// ----------------------------------------------------------------------------

func logOSInfo() {
	cmd := exec.Command("ps", "-eo", "nlwp,cmd")
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println("DEBUG OS:", string(stdout))
}

// ----------------------------------------------------------------------------

func logStats() {
	cpus := runtime.NumCPU()
	goMaxProc := runtime.GOMAXPROCS(0)
	goRoutines := runtime.NumGoroutine()
	cgoCalls := runtime.NumCgoCall()
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	var gcStats debug.GCStats
	debug.ReadGCStats(&gcStats)

	// fmt.Println("---------------------------------------------------------------")
	// fmt.Println("Time:", time.Now())
	// fmt.Println("CPUs:", cpus)
	// fmt.Println("GOMAXPROC:", goMaxProc)
	// fmt.Println("Go routines:", goRoutines)
	// fmt.Println("CGO calls:", cgoCalls)
	// fmt.Println("Num GC:", memStats.NumGC)
	// fmt.Println("GCSys:", memStats.GCSys)
	// fmt.Println("GC pause total:", gcStats.PauseTotal)
	// fmt.Println("LastGC:", gcStats.LastGC)
	// fmt.Println("HeapAlloc:", memStats.HeapAlloc)
	// fmt.Println("NextGC:", memStats.NextGC)
	// fmt.Println("CPU fraction used by GC:", memStats.GCCPUFraction)

	fmt.Println("---------------------------------------------------------------")
	printCSV(">>>", "Time", "CPUs", "GOMAXPROC", "Go routines", "CGO calls", "Num GC", "GC pause total", "LastGC", "TotalAlloc", "HeapAlloc", "NextGC", "GCSys", "HeapSys", "StackSys", "Sys - total OS bytes", "CPU fraction used by GC")
	printCSV(">>>", time.Now(), cpus, goMaxProc, goRoutines, cgoCalls, memStats.NumGC, gcStats.PauseTotal, gcStats.LastGC, memStats.TotalAlloc, memStats.HeapAlloc, memStats.NextGC, memStats.GCSys, memStats.HeapSys, memStats.StackSys, memStats.Sys, memStats.GCCPUFraction)
}

// ----------------------------------------------------------------------------

func printCSV(fields ...any) {
	for _, field := range fields {
		fmt.Print(field, ",")
	}
	fmt.Println("")
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
