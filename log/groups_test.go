package log

import (
	"bytes"
	"github.com/alt4dev/go/service"
	"os"
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

func TestGrouping(t *testing.T){
	defer CloseGroup()
	// Calling groupId without Initializing a group will create a different group on every call
	if getGroupId() == getGroupId() {
		t.Error("If group is not initialized, call to groupId should return a unique value every time")
	}

	InitGroup(nil, nil)

	groupId, ok := logGroups[getRoutineId()]
	if !ok {
		t.Error("Unable to find initialized group")
	}
	if getGroupId() != groupId || getGroupId() != getGroupId() {
		t.Error("After group initialization, all calls to getGroupId should return the same groupId")
	}

	// Capture logs from here
	var buffer bytes.Buffer
	service.EmitWarning.SetOutput(&buffer)
	defer func(){
		service.EmitWarning.SetOutput(os.Stderr)
	}()

	// Reinitializing a group should delete older group and start a new one
	InitGroup(nil, nil)
	// Confirm that we warned the user
	warning := "Unclosed log group detected. Call `defer CloseGroup()` after `InitGroup()` to avoid memory leaks"
	if !strings.Contains(buffer.String(), warning){
		t.Error("Replacing a log group should warn the user of the previously unclosed log group")
	}

	if groupId == getGroupId() {
		t.Error("Initializing a log group should use a new groupId")
	}
}

func TestGroupClosing(t *testing.T){
	InitGroup(nil, nil)
	routine := getRoutineId()
	if _, ok := logGroups[routine]; !ok {
		t.Error("Unable to find initialized group")
	}

	CloseGroup()
	// Confirm that we deleted the group from memory
	if _, ok := logGroups[routine]; ok {
		t.Error("Group found after closing")
	}
}