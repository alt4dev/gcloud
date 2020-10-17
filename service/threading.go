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

var threads map[string]string
var waitGroups map[string] *sync.WaitGroup

func init()  {
	threads = map[string]string{}
	waitGroups = map[string]*sync.WaitGroup{}
}

func getRoutineId() string {
	traces := strings.Split(string(debug.Stack()), "\n")
	return traces[0]
}

func getThreadId() string {
	routineId := getRoutineId()
	if val, ok := threads[routineId]; ok {
		return val
	}
	// Return a uuid if not grouped
	return uuid.New().String()
}

func initGroup() {
	routineId := getRoutineId()
	if _, ok := threads[routineId]; ok {
		emitWarning.Println("Unclosed log group detected. Call `defer group.Close()` after initializing group to avoid memory leaks. Better yet do `defer Group(title, claims).Close()`")
		delete(threads, routineId)
	}
	threads[routineId] = getThreadId()
}

func CloseGroup() {
	// Before closing a group. Wait for all logs to finish writing.
	WaitGroup().Wait()

	routineId := getRoutineId()
	if _, ok := threads[routineId]; ok {
		delete(threads, routineId)
	}
}

// Provide wait groups per go routine ID. Closing a group will wait for all write ops to finish.
func WaitGroup() *sync.WaitGroup {
	routineId := getRoutineId()
	wg, ok := waitGroups[routineId]
	if ok {
		return wg
	}
	wg = &sync.WaitGroup{}
	waitGroups[routineId] = wg
	return wg
}

