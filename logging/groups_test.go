package logging

import (
	"runtime/debug"
	"strings"
	"testing"
)


func TestRoutineId(t *testing.T) {
	// We take routine ID from the first line of the trace
	traces := strings.Split(string(debug.Stack()), "\n")
	routineId := getRoutineId()
	if traces[0] != routineId {
		t.Error("Test found different routines, expected:\n", traces[0], "\nbut got:\n", routineId)
	}

	if routineId != getRoutineId() {
		t.Error("Continuous calls to getRoutineId() should return the same Id")
	}
}
