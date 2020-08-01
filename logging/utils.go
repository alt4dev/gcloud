package logging

import (
	"github.com/google/uuid"
	"log"
	"os"
	"runtime/debug"
	"strings"
)

var (
	alt4error *log.Logger
	alt4warning *log.Logger
	instanceId string
	logGroups map[string]string
)

func init()  {
	alt4error = log.New(os.Stderr, "[alt4](if seeing this please contact: critical@alt4.dev) ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	alt4warning = log.New(os.Stderr, "[alt4] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
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


