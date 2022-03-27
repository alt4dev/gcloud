package service

import (
	"cloud.google.com/go/logging"
	"context"
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/google/uuid"
	loggingProto "google.golang.org/genproto/googleapis/appengine/logging/v1"
	ltype "google.golang.org/genproto/googleapis/logging/type"
	logpb "google.golang.org/genproto/googleapis/logging/v2"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/durationpb"
	"google.golang.org/protobuf/types/known/timestamppb"
	"net/http"
	"runtime"
	"strings"
	"time"
)

var levels []logging.Severity

func init() {
	levels = []logging.Severity{logging.Default, logging.Info, logging.Debug, logging.Warning, logging.Error, logging.Critical, logging.Alert}
}

func numLogLevel(level logging.Severity) int {
	for k, v := range levels {
		if level == v {
			return k
		}
	}
	return -1
}

func makeTrace(id string) string {
	return fmt.Sprintf("projects/%s/traces/%s", project, id)
}

// Log Creates a log entry and writes it to Google cloud logging in the background.
// This function should not be called directly and should instead be used from helper functions under the `log` package.
func Log(callDepth int, message string, level logging.Severity, logTime time.Time) {
	// Get the parent file and function of the caller
	pc, file, line, _ := runtime.Caller(callDepth)
	function := runtime.FuncForPC(pc).Name()

	if options.Mode == ModeTesting || options.Mode == ModeDebug {
		emitLog(message, nil, level, logTime, file, line)
	}

	// Early return for testing and silent modes
	if options.Mode == ModeTesting || options.Mode == ModeSilent {
		return
	}

	details := getThreadDetails()

	entry := logging.Entry{
		Timestamp: logTime,
		Severity:  level,
		Payload:   message,
		InsertID:  uuid.New().String(),
		Resource:  options.Resource,
		SourceLocation: &logpb.LogEntrySourceLocation{
			File:     file,
			Line:     int64(line),
			Function: function,
		},
	}

	// We can't group logs without a request
	if details == nil || details.httpRequest == nil {
		// We can however add a trace and operation
		if details != nil {
			entry.Trace = makeTrace(details.id)
			entry.Operation = &logpb.LogEntryOperation{
				Id:       details.id,
				Producer: "github.com/alt4dev/gcloud",
			}
		}
		client.Logger(getThreadId()).Log(entry)
		return
	}

	// Lock detail momentarily
	details.lock.Lock()
	defer details.lock.Unlock()

	// Add a new line if the message is set
	if message != "" {
		logLine := loggingProto.LogLine{
			LogMessage: message,
			Severity:   ltype.LogSeverity(level),
			SourceLocation: &loggingProto.SourceLocation{
				File:         file,
				FunctionName: function,
				Line:         int64(line),
			},
			Time: timestamppb.New(logTime),
		}

		details.lines = append(details.lines, &logLine)
	}

	numLevel := numLogLevel(level)

	if numLevel > details.level {
		details.level = numLevel
		// Update details with the highest loglevel
	}
}

func readUserIP(r *http.Request) string {
	IPAddress := r.Header.Get("X-Real-Ip")
	if IPAddress == "" {
		IPAddress = r.Header.Get("X-Forwarded-For")
	}
	if IPAddress == "" {
		IPAddress = r.RemoteAddr
	}
	return IPAddress
}

func writeGroupLog(details *threadDetails, status int, labels map[string]string) {
	// Return if we don't have a request closing the log.

	request := details.httpRequest

	endTime := LogTime()
	entry := logging.Entry{
		Timestamp: details.startTime,
		Severity:  levels[details.level],
		Labels:    labels,
		InsertID:  uuid.New().String(),
		HTTPRequest: &logging.HTTPRequest{
			Request:     request,
			RequestSize: request.ContentLength,
			Status:      status,
			Latency:     endTime.Sub(details.startTime),
			LocalIP:     "",
			RemoteIP:    readUserIP(request),
		},
		Operation: &logpb.LogEntryOperation{
			Id:       details.id,
			Producer: "github.com/alt4dev/gcloud",
			First:    true,
			Last:     true,
		},
		Resource: options.Resource,
		Trace:    makeTrace(details.id),
	}

	pbEntry, err := logging.ToLogEntry(entry, project)
	if err != nil {
		client.OnError(err)
		return
	}

	protocMsg := &loggingProto.RequestLog{
		RequestId:     details.id,
		Ip:            request.RemoteAddr,
		StartTime:     timestamppb.New(details.startTime),
		EndTime:       timestamppb.Now(),
		Latency:       durationpb.New(entry.HTTPRequest.Latency),
		Method:        request.Method,
		Resource:      request.URL.RawPath,
		Status:        int32(status),
		ResponseSize:  0,
		Referrer:      request.Referer(),
		UserAgent:     request.UserAgent(),
		Host:          request.Host,
		TaskQueueName: request.Header.Get("X-AppEngine-QueueName"),
		TaskName:      request.Header.Get("X-AppEngine-TaskName"),
		Finished:      true,
		First:         true,
		Line:          details.lines,
		TraceId:       entry.Trace,
		TraceSampled:  false,
		// TODO: Support adding source reference e.g. Github
		SourceReference: nil,
	}

	protocBytes, err := proto.Marshal(protocMsg)
	if err != nil {
		client.OnError(err)
		return
	}

	// Set Payload to the RequestLog we just created
	pbEntry.Payload = &logpb.LogEntry_ProtoPayload{ProtoPayload: &anypb.Any{
		TypeUrl: "type.googleapis.com/google.appengine.logging.v1.RequestLog",
		Value:   protocBytes,
	}}

	req := &logpb.WriteLogEntriesRequest{
		LogName:  details.id,
		Resource: options.Resource,
		Labels:   nil,
		Entries:  []*logpb.LogEntry{pbEntry},
	}

	retriedWrite(req, 0)
}

func retriedWrite(req *logpb.WriteLogEntriesRequest, retries int) {
	if retries >= 10 {
		emitError.Print("Unable to write log entry to google cloud logging")
		emitError.Print(req.Entries[0])
		return
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*60)
	defer cancel()

	_, err := pbClient.WriteLogEntries(ctx, req)
	if err != nil {
		client.OnError(err)
		retriedWrite(req, retries+1)
	}
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
