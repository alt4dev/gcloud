package service

import (
	"context"
	"errors"
	"github.com/alt4dev/protobuff/proto"
	"runtime"
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
func Log(threadInit bool, message string, claims []*proto.Claim, level uint8) *LogResult {
	if threadInit {
		initGroup()
	}
	logTime :=time.Now()
	// Get the parent file and function of the caller
	pc, file, line, _ := runtime.Caller(2)
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
	//defer result.wg.Done()
	if getClient() == nil {
		result.error = errors.New("error connecting to remote server")
		result.wg.Done()
		return
	}
	result.result, result.error = (*client).Log(context.Background(), msg)
	result.wg.Done()
}