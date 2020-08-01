package logging

import (
	"github.com/google/uuid"
	"runtime/debug"
	"strings"
)

var (
	instanceId string
	logGroups map[string]string
)

func init()  {
	instanceId = uuid.New().String()
	logGroups = make(map[string]string, 0)
}


func getRoutineId() string {
	traces := strings.Split(string(debug.Stack()), "\n")
	return traces[0]
}

func getGroupId() string {
	routineId := getRoutineId()
	if val, ok := logGroups[routineId]; ok {
		return val
	}
	// Return a uuid if not grouped
	return uuid.New().String()
}

func InitGroup(groupName *string, claims Claim){
	routineId := getRoutineId()
	if _, ok := logGroups[routineId]; ok {
		alt4warning.Println("Unclosed log group detected. Call `defer CloseGroup()` after `InitGroup()` to avoid memory leaks")
		delete(logGroups, routineId)
	}
	logGroups[routineId] = uuid.New().String()
}

func CloseGroup() {
	routineId := getRoutineId()
	if _, ok := logGroups[routineId]; ok {
		delete(logGroups, routineId)
	}
}


