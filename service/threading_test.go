package service

import (
	"os"
	"sync"
	"testing"
)

type routineResult struct {
	wg        *sync.WaitGroup
	t         *testing.T
	routineId string
}

func runInThread(result *routineResult) {
	result.routineId = getRoutineId()
	// Confirm that within the same routine the ID is the same
	if result.routineId != getRoutineId() {
		result.t.Error("Expected routine Id to remain the same in the same routine")
	}
	result.wg.Done()
}

func TestGetRoutineId(t *testing.T) {
	// Test that starting different go routines will return different routine IDs
	waitGroup := &sync.WaitGroup{}
	waitGroup.Add(3)

	threadA := routineResult{
		wg: waitGroup,
	}
	threadB := routineResult{
		wg: waitGroup,
	}
	threadC := routineResult{
		wg: waitGroup,
	}
	go runInThread(&threadA)
	go runInThread(&threadB)
	go runInThread(&threadC)

	// Wait for all threads to finish
	waitGroup.Wait()

	// Confirm that none of the routines match each other
	if threadA.routineId == threadB.routineId || threadB.routineId == threadC.routineId || threadC.routineId == threadA.routineId {
		t.Error("Go routines returned the same routine ID.")
		return
	}
}

func TestGrouping(t *testing.T) {
	// Test that before grouping thread ID keeps changing
	if getThreadId() == getThreadId() {
		t.Error("Thread IDs should not match for ungrouped routines")
		return
	}

	// Start a group
	initGroup()

	currentThread := getThreadId()
	// Confirm that with grouping, the thread ID doesn't change
	if currentThread != getThreadId() {
		t.Error("Thread IDs are expected to be the same once a group is initialized")
		return
	}

	f, _ := os.Open(os.DevNull)
	EmitWarning.SetOutput(f) // suppress warning log
	// Confirm that initializing a group without closing replaces the older one
	initGroup()
	EmitWarning.SetOutput(options.Writer) // restore
	newThreadId := getThreadId()
	if newThreadId == currentThread {
		t.Error("Thread ID's should change after initializing an existing group")
	}

	if _, ok := threads[getRoutineId()]; !ok {
		t.Error("In an initialized group, routine id should be part of threads list")
	}

	CloseGroup()
	if _, ok := threads[getRoutineId()]; ok {
		t.Error("Routine id should be deleted from threads list after closing a thread")
	}
}
