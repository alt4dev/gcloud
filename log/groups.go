package log

import (
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
)

// Group start a log group for the goroutine that calls this function. A group should be closed after. Use: `defer Group(title, claims).Close()`
// Go routines are used to write logs to alt4 without blocking. You can use the module `sync/log` to write synchronously or call the function `Result` which is going to block if the operation is not done.
func Group(title string, claims Claims) GroupResult {
	return GroupResult{
		logResult: service.Log(2,true, title, claims.parse(), 1),
	}
}

// GroupResult Object returned by creating a new log group/thread.
type GroupResult struct {
	logResult *service.LogResult
}

func (result GroupResult) Result() (*proto.Result, error) {
	return result.logResult.Result()
}

func (result GroupResult) Close() {
	service.CloseGroup()
}
