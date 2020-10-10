package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/alt4dev/protobuff/proto"
	"runtime"
	"strings"
	"sync"
	"time"
)

// LogResult Object returned when you create a log entry.
type LogResult struct {
	wg *sync.WaitGroup
	result *proto.Result
	error error
}

// Result Returns actual result from alt4. This will block and wait for the result if not done
func (result *LogResult) Result() (*proto.Result, error) {
	result.wg.Wait()
	return result.result, result.error
}

// Log Creates a log entry and writes it to alt4 in the background.
// This function should not be called directly and should instead be used from helper functions under the `log` package.
func Log(calldepth int, threadInit bool, message string, claims []*proto.Claim, level uint8) *LogResult {
	if threadInit {
		initGroup()
	}
	logTime :=time.Now()
	// Get the parent file and function of the caller
	pc, file, line, _ := runtime.Caller(calldepth)
	function := runtime.FuncForPC(pc).Name()
	msg := proto.Message{
		ThreadId:   getThreadId(),
		SinkId:     options.Sink,
		ThreadInit: threadInit,
		Message:    message,
		Claims:     claims,
		FileName:   file,
		LineNo:     uint32(line),
		Function: 	function,
		Level:      uint32(level),
		Timestamp:  uint64(logTime.UnixNano()),
		AuthToken:  options.AuthToken,
	}
	result := LogResult{
		wg: &sync.WaitGroup{},
	}
	result.wg.Add(1)
	go writeToAlt4(&msg, &result)
	return &result
}

func writeToAlt4(msg *proto.Message, result *LogResult){
	if getClient() == nil {
		result.error = errors.New("error connecting to remote server")
		emitLog(msg, result.error)
		result.wg.Done()
		return
	}
	if options.Mode != "testing" {
		result.result, result.error = (*client).Log(context.Background(), msg)
	}
	shouldEmit := options.Mode == "debug" || options.Mode == "testing"
	if (result.result != nil && !result.result.Acknowledged) || shouldEmit || result.error != nil {
		if result.result != nil && !result.result.Acknowledged {
			result.error = errors.New(result.result.Message)
		}
		emitLog(msg, result.error)
	}
	result.wg.Done()
}

func emitLog(msg *proto.Message, err error) {
	if err != nil {
		EmitError.Println(err)
	}
	timeString := time.Unix(0, int64(msg.Timestamp)).Format("2021-12-15 13:45:45.000")
	message := fmt.Sprintf("%s %s %s:%d %s", levelString[uint8(msg.Level)], timeString, msg.FileName, msg.LineNo, msg.Message)
	lines := []string{message}
	for _, claim := range msg.Claims {
		lines = append(lines, fmt.Sprintf("\tclaim.%s: '%s'", claim.Name, claim.Value))
	}
	Emit.Println(strings.Join(lines, "\n"))
}