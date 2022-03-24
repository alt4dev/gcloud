package log

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/alt4dev/gcloud/service"
	"net/http"
	"net/url"
	"runtime/debug"
	"time"
)

// GroupResult Object returned by creating a new log group/thread.
type GroupResult struct {
	labels    *Labels
	startTime time.Time
	status    int
	method    string
	url       string
	request   *http.Request
}

// Close will mark the end of a thread closing the log group.
// If arguments are provided to the close function, they'll be logged.
// This can be useful for determining the latency of a request.
// If there were unfinished writes to alt4 during this thread.
// This method will wait for the writes to finish
// Close also logs any panic but doesn't recover.
func (result *GroupResult) Close() {
	defer service.CloseGroup()
	var labels map[string]string = nil
	if result.labels != nil {
		labels = result.labels.parse()
	}

	r := recover()
	// Recover any panic, just to log
	if r != nil {
		service.Log(2, fmt.Sprint(r), nil, labels, logging.Critical, service.LogTime())
		// Log stack trace
		service.Log(2, fmt.Sprint(string(debug.Stack())), nil, labels, logging.Critical, service.LogTime())
	}

	// Parent log should have the last timestamp in the group
	logTime := service.LogTime()

	latency := logTime.Sub(result.startTime)

	httpRequest := result.request
	if httpRequest == nil {
		_url, _ := url.Parse(result.url)
		httpRequest = &http.Request{
			Method: result.method,
			URL:    _url,
		}
	}

	// Complete the request
	request := &logging.HTTPRequest{
		Request: httpRequest,
		Status:  result.status,
		Latency: latency,
	}

	message := fmt.Sprintf("[%s] %s", result.method, result.url)

	service.Log(2, message, request, labels, logging.Default, logTime)
	// Anakin, continue panakin
	if r != nil {
		panic(r)
	}
}

func (result *GroupResult) SetStatus(httpStatus int) {
	result.status = httpStatus
}

func (result *GroupResult) SetRequest(request *http.Request) {
	result.request = request
}
