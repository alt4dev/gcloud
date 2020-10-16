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
	R   *proto.Result
	wg *sync.WaitGroup
	Err error
}

// Result Returns actual R from alt4. This will block and wait for the R if not done
func (result *LogResult) Result() (*proto.Result, error) {
	result.wg.Wait()
	return result.R, result.Err
}

// RemoteWriter an interface for functions called when writing to alt4.
// You can implement this function to mock writes to alt4 for better testing of your system.
type RemoteWriter interface {
	// Write function will be called with the Message to be sent to alt4 and an empty R to fill once done
	Write(msg *proto.Message, result *LogResult)
}

type writer struct{}

func (w writer) Write(msg *proto.Message, result *LogResult) {
	if getClient() == nil {
		result.Err = errors.New("error connecting to remote server")
		emitLog(msg, result.Err)
		return
	}
	if options.Mode != "testing" {
		result.R, result.Err = (*client).Log(context.Background(), msg)
	}
	shouldEmit := options.Mode == "debug" || options.Mode == "testing"
	if (result.R != nil && !result.R.Acknowledged) || shouldEmit || result.Err != nil {
		if result.R != nil && !result.R.Acknowledged {
			result.Err = errors.New(result.R.Message)
		}
		emitLog(msg, result.Err)
	}
}

// Alt4RemoteWriter For testing purposes, implement your own RemoteWriter and equate it to this variable
var Alt4RemoteWriter RemoteWriter = writer{}

// Log Creates a log entry and writes it to alt4 in the background.
// This function should not be called directly and should instead be used from helper functions under the `log` package.
func Log(calldepth int, threadInit bool, message string, claims []*proto.Claim, level uint8) *LogResult {
	if threadInit {
		initGroup()
	}
	logTime := time.Now()
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
		Function:   function,
		Level:      uint32(level),
		Timestamp:  uint64(logTime.UnixNano()),
		AuthToken:  options.AuthToken,
	}
	result := LogResult{
		wg: WaitGroup(),
	}
	WaitGroup().Add(1)
	go writerHelper(&msg, &result)
	return &result
}

func writerHelper(msg *proto.Message, result *LogResult) {
	defer result.wg.Done()
	Alt4RemoteWriter.Write(msg, result)
}

func emitLog(msg *proto.Message, err error) {
	if err != nil {
		EmitError.Println(err)
	}
	timeString := time.Unix(0, int64(msg.Timestamp)).Format("2006-01-02T15:04:05.000Z")
	message := fmt.Sprintf("[alt4 %s] %s %s:%d %s", levelString[uint8(msg.Level)], timeString, msg.FileName, msg.LineNo, msg.Message)
	lines := []string{message}
	for _, claim := range msg.Claims {
		lines = append(lines, fmt.Sprintf("\tclaim.%s: '%s'", claim.Name, claim.Value))
	}
	Emit.Println(strings.Join(lines, "\n"))
}
