package service

import (
	"cloud.google.com/go/logging"
	"github.com/google/uuid"
	loggingProto "google.golang.org/genproto/googleapis/appengine/logging/v1"
	"net/http"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

/*
Golang doesn't support threading in the traditional way.
This file has functions that try at best effort to identify the beginning and the end of a
goroutine and put those logs in the same group
*/

var threads sync.Map

func init() {
	threads = sync.Map{}
}

type threadDetails struct {
	id          string
	level       int
	lines       []*loggingProto.LogLine
	logger      *logging.Logger
	lock        *sync.Mutex
	startTime   time.Time
	httpRequest *http.Request
}

func getRoutineId() string {
	traces := strings.Split(string(debug.Stack()), "\n")
	return traces[0]
}

func getThreadId() string {
	routineId := getRoutineId()
	if val, ok := threads.Load(routineId); ok {
		return val.(string)
	}
	// Return a new uuid if not grouped
	return uuid.New().String()
}

func InitGroup(httpRequest *http.Request) time.Time {
	routineId := getRoutineId()
	if detailsInterface, ok := threads.Load(routineId); ok {
		details := detailsInterface.(*threadDetails)
		_ = details.logger.Flush()
		threads.Delete(routineId)
		emitWarning.Println("Unclosed log group detected. Call `defer group.Close()` after initializing group to avoid memory leaks. Better yet do `defer Group(request, nil).Close()`")
	}

	if httpRequest == nil {
		emitWarning.Println("Group with a `nil` request. These logs won't be grouped correctly. However each group will have a unique trace which you can use to filter your logs by and will be a member of the same operation.")
	}
	threadId := getThreadId()
	startTime := LogTime()
	threads.Store(routineId, &threadDetails{
		id:          threadId,
		level:       0,
		lines:       make([]*loggingProto.LogLine, 0),
		logger:      client.Logger(threadId),
		lock:        new(sync.Mutex),
		startTime:   startTime,
		httpRequest: httpRequest,
	})

	return startTime
}

func getThreadDetails() *threadDetails {
	routineId := getRoutineId()
	if detailsInterface, ok := threads.Load(routineId); ok {
		return detailsInterface.(*threadDetails)
	}
	return nil
}

func CloseGroup(status int, labels map[string]string) {
	routineId := getRoutineId()
	if detailsInterface, ok := threads.Load(routineId); ok {
		// Delete the thread no matter what happens
		defer threads.Delete(routineId)
		details := detailsInterface.(*threadDetails)
		if details.httpRequest != nil {
			// Logs are written as a group
			writeGroupLog(details, status, labels)
		}
		// Before closing a group. Wait for all logs to finish writing.
		_ = details.logger.Flush()
	}
}
