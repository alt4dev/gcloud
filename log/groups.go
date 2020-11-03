package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
	"time"
)

// GroupResult Object returned by creating a new log group/thread.
type GroupResult struct {
	logResult *service.LogResult
	claims *Claims
}

// Return the result of the actual log event
func (result GroupResult) Result() (*proto.Result, error) {
	return result.logResult.Result()
}

// Close will mark the end of a thread closing the log group.
// If arguments are provided to the close function, they'll be logged.
// This can be useful for determining the latency of a request.
// If there were unfinished writes to alt4 during this thread.
// This method will wait for the writes to finish
func (result GroupResult) Close(v ...interface{}) {
	t := time.Now()
	if len(v) > 0{
		var claims []*proto.Claim = nil
		if result.claims != nil {
			claims = result.claims.parse()
		}
		message := fmt.Sprint(v...)
		service.Log(2, false, message, claims, proto.Log_NONE, t)
	}
	service.CloseGroup()
}
