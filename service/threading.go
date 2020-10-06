package service

import (
	"github.com/google/uuid"
	"runtime/debug"
	"strings"
)

/*
Golang doesn't support threading in the traditional way.
This file has functions that try at best effort to identify the beginning and the end of a
goroutine and put those logs in the same group
*/

var threads map[string]string

func init()  {
	threads = make(map[string]string, 0)
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
		EmitWarning.Println("Unclosed log group detected. Call `defer group.Close()` after initializing group to avoid memory leaks. Better yet do `defer Group(title, claims).Close()`")
		delete(threads, routineId)
	}
	threads[routineId] = getThreadId()
}

func CloseGroup() {
	routineId := getRoutineId()
	if _, ok := threads[routineId]; ok {
		delete(threads, routineId)
	}
}


