package log

import (
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
)

// GroupResult Object returned by creating a new log group/thread.
type GroupResult struct {
	logResult *service.LogResult
}

// Return the result of the actual log event
func (result GroupResult) Result() (*proto.Result, error) {
	return result.logResult.Result()
}

// Close will mark the end of a thread closing the log group.
// If there were unfinished writes to alt4 during this thread.
// This method will wait for the writes to finish
func (result GroupResult) Close() {
	service.CloseGroup()
}
