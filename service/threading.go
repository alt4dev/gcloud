package service

import (
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
var waitGroups sync.Map

func init()  {
	threads = sync.Map{}
	waitGroups = sync.Map{}
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
	// Return a uuid if not grouped
	return uuid.New().String()
}

func initGroup() {
	routineId := getRoutineId()
	if _, loaded := threads.LoadAndDelete(routineId); loaded {
		emitWarning.Println("Unclosed log group detected. Call `defer group.Close()` after initializing group to avoid memory leaks. Better yet do `defer Group(title, claims).Close()`")
	}
	threads.Store(routineId, getThreadId())
}

func CloseGroup() {
	// Before closing a group. Wait for all logs to finish writing.
	WaitGroup().Wait()

	routineId := getRoutineId()
	if _, loaded := threads.LoadAndDelete(routineId); loaded {
		waitGroups.Delete(routineId)
	}
}

// Provide wait groups per go routine ID. Closing a group will wait for all write ops to finish.
func WaitGroup() *sync.WaitGroup {
	routineId := getRoutineId()
	wg, ok := waitGroups.Load(routineId)
	if ok {
		return wg.(*sync.WaitGroup)
	}
	wg = &sync.WaitGroup{}
	waitGroups.Store(routineId, wg)
	return wg.(*sync.WaitGroup)
}

