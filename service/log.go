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

// Log Creates a log entry and writes it to alt4 in the background.
// This function should not be called directly and should instead be used from helper functions under the `log` package.
func Log(callDepth int, asGroup bool, message string, request *logging.HTTPRequest, labels map[string]string, level logging.Severity, logTime time.Time) {
	// Don't use resources grouping for testing modes
	if asGroup && (options.Mode != ModeTesting && options.Mode != ModeSilent) {
		initGroup()
	}

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
		Labels:      nil,
		InsertID:    uuid.New().String(),
		HTTPRequest: request,
		Resource:    nil,
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

	// Set trace ID to the thread ID
	entry.Trace = details.id

	logger := details.child

	if request != nil {
		logger = details.parent
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
