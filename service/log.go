package service

import (
	"context"
	"errors"
	"github.com/alt4dev/protobuff/proto"
	"sync"
	"time"
)

// LogResult Object returned when you create a log entry.
type LogResult struct {
	wg sync.WaitGroup
	result *proto.Result
	error error
}

// Result Returns actual result from alt4. This will block and wait for the result if not done
func (result LogResult) Result() (*proto.Result, error) {
	result.wg.Wait()
	return result.result, result.error
}

func Log(threadId string, sinkId string, threadInit bool, message string, claims []*proto.Claim, file string, line uint32, function string, level uint8) LogResult {
	msg := proto.Message{
		ThreadId:   threadId,
		SinkId:     sinkId,
		ThreadInit: threadInit,
		Message:    message,
		Claims:     claims,
		FileName:   file,
		LineNo:     line,
		Function: 	function,
		Level:      uint32(level),
		Timestamp:  uint64(time.Now().UnixNano()),
		AuthToken:  options.AuthToken,
	}

	result := LogResult{}
	result.wg.Add(1)
	go func() {
		defer result.wg.Done()
		if getClient() == nil {
			result.error = errors.New("error connecting to remote server")
			return
		}
		result.result, result.error = (*client).Log(context.Background(), &msg)
	}()
	return result
}
