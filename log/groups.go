package log

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/alt4dev/gcloud/service"
	"net/http"
	"runtime/debug"
	"time"
)

// Labels are fields that will can be associated to your group log entry.
// They can be used to filter and better identify your logs.
type Labels map[string]interface{}

func (labels Labels) parse() map[string]string {
	protoLabels := map[string]string{}
	for key, i := range labels {
		var claimValue string
		switch i.(type) {
		case bool:
			claimValue = fmt.Sprint(i.(bool))
		case time.Time:
			claimValue = i.(time.Time).Format("2006-01-02T15:04:05.000Z")
		default:
			claimValue = fmt.Sprint(i)
		}
		protoLabels[key] = claimValue
	}
	return protoLabels
}

// GroupResult Object returned by creating a new log group/thread.
type GroupResult struct {
	labels Labels
	status int
}

// Close will mark the end of a thread closing the log group.
// If arguments are provided to the close function, they'll be logged.
// This can be useful for determining the latency of a request.
// If there were unfinished writes to alt4 during this thread.
// This method will wait for the writes to finish
// Close also logs any panic but doesn't recover.
func (result *GroupResult) Close() {
	var labels map[string]string = nil
	if result.labels != nil {
		labels = result.labels.parse()
	}

	defer service.CloseGroup(result.status, labels)

	r := recover()
	// Recover any panic for recording, then continue
	if r != nil {
		service.Log(2, fmt.Sprint(r), logging.Critical, service.LogTime())
		// Log stack trace
		service.Log(2, fmt.Sprint(string(debug.Stack())), logging.Critical, service.LogTime())

		// Anakin, continue panakin
		panic(r)
	}
}

func (result *GroupResult) SetStatus(httpStatus int) *GroupResult {
	result.status = httpStatus
	return result
}

func (result *GroupResult) SetLabel(key string, value interface{}) *GroupResult {
	if result.labels == nil {
		result.labels = Labels{}
	}

	result.labels[key] = value
	return result
}

// Group start a log group for the goroutine that calls this function.
// A group should be closed after. Use: `defer Group(request, nil).Close()`
func Group(request *http.Request, labels Labels) *GroupResult {
	service.InitGroup(request)
	return &GroupResult{
		labels: labels,
		status: 0,
	}
}
