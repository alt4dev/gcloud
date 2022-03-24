package service

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/google/uuid"
	logging2 "google.golang.org/genproto/googleapis/logging/v2"
	"runtime"
	"strings"
	"time"
)

func numLogLevel(level logging.Severity, levels []logging.Severity) int {
	for k, v := range levels {
		if level == v {
			return k
		}
	}
	return -1
}

// Log Creates a log entry and writes it to Google cloud logging in the background.
// This function should not be called directly and should instead be used from helper functions under the `log` package.
func Log(callDepth int, message string, request *logging.HTTPRequest, labels map[string]string, level logging.Severity, logTime time.Time) {
	// Get the parent file and function of the caller
	pc, file, line, _ := runtime.Caller(callDepth)
	function := runtime.FuncForPC(pc).Name()

	if options.Mode == ModeTesting || options.Mode == ModeDebug {
		emitLog(message, labels, level, logTime, file, line)
	}

	// Early return for testing and silent modes
	if options.Mode == ModeTesting || options.Mode == ModeSilent {
		return
	}

	details := getThreadDetails()

	entry := logging.Entry{
		Timestamp:   logTime,
		Severity:    level,
		Payload:     message,
		Labels:      labels,
		InsertID:    uuid.New().String(),
		HTTPRequest: request,
		Resource:    options.Resource,
		SourceLocation: &logging2.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: function,
		},
	}

	if details == nil {
		client.Logger(getThreadId()).Log(entry)
		return
	}

	levels := []logging.Severity{logging.Default, logging.Info, logging.Debug, logging.Warning, logging.Error, logging.Critical, logging.Alert}
	numLevel := numLogLevel(level, levels)

	if numLevel > details.level {
		details.level = numLevel
		// Update details with the highest loglevel
		threads.Store(getRoutineId(), details)
	}

	// Set trace ID to the thread ID
	entry.Trace = details.id

	logger := details.child

	if request != nil {
		logger = details.parent
		// Set the parent log level to the highest level seen
		entry.Severity = levels[details.level]
	}

	logger.Log(entry)
}

func emitLog(message string, labels map[string]string, level logging.Severity, logTime time.Time, file string, line int) {
	timeString := logTime.Format("2006-01-02T15:04:05.000Z")
	message = fmt.Sprintf("[alt4 %s] %s %s:%d %s", level.String(), timeString, file, line, message)
	lines := []string{message}
	for key, value := range labels {
		lines = append(lines, fmt.Sprintf("\tlabels.%s: '%s'", key, value))
	}
	emit.Println(strings.Join(lines, "\n"))
}
