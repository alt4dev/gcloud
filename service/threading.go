package service

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/google/uuid"
	"runtime/debug"
	"strings"
	"sync"
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
	id     string
	level  int
	parent *logging.Logger
	child  *logging.Logger
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

func InitGroup() {
	routineId := getRoutineId()
	if detailsInterface, ok := threads.Load(routineId); ok {
		details := detailsInterface.(*threadDetails)
		_ = details.child.Flush()
		_ = details.parent.Flush()
		threads.Delete(routineId)
		emitWarning.Println("Unclosed log group detected. Call `defer group.Close()` after initializing group to avoid memory leaks. Better yet do `defer Group(url, method).Close()`")
	}
	threadId := getThreadId()
	threads.Store(routineId, &threadDetails{
		id:     threadId,
		level:  0,
		parent: client.Logger(fmt.Sprintf("parent-%s", threadId)),
		child:  client.Logger(threadId),
	})
}

func getThreadDetails() *threadDetails {
	routineId := getRoutineId()
	if detailsInterface, ok := threads.Load(routineId); ok {
		return detailsInterface.(*threadDetails)
	}
	return nil
}

func CloseGroup() {
	routineId := getRoutineId()
	if detailsInterface, ok := threads.Load(routineId); ok {
		details := detailsInterface.(*threadDetails)
		// Before closing a group. Wait for all logs to finish writing.
		_ = details.child.Flush()
		_ = details.parent.Flush()
		threads.Delete(routineId)
	}
}
