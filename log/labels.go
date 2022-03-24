package log

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/alt4dev/gcloud/service"
	"time"
)

// Labels are fields that will can be associated to your log entry.
// They can be used to filter and better identify your logs.
type Labels map[string]interface{}

// Group start a log group for the goroutine that calls this function.
// A group should be closed after. Use: `defer Labels{...}.Group(...).Close()`
func (labels Labels) Group(url string, method string) *GroupResult {

	return &GroupResult{
		labels:    &labels,
		startTime: service.LogTime(),
		status:    0,
		method:    method,
		url:       url,
		request:   nil,
	}
}

// Print send labels and a log message to Google cloud logging
func (labels Labels) Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, labels.parse(), logging.Default, service.LogTime())
}

// Printf send labels and a log message to Google cloud logging
func (labels Labels) Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, labels.parse(), logging.Default, service.LogTime())
}

// Println send labels and a log message to Google cloud logging
func (labels Labels) Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, message, nil, labels.parse(), logging.Default, service.LogTime())
}

// Info send labels and a logging.Info message to Google cloud logging
func (labels Labels) Info(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, labels.parse(), logging.Info, service.LogTime())
}

// Infof send labels and a logging.Info message to Google cloud logging
func (labels Labels) Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, labels.parse(), logging.Info, service.LogTime())
}

// Debug send labels and a logging.Debug message to Google cloud logging
func (labels Labels) Debug(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, labels.parse(), logging.Debug, service.LogTime())
}

// Debugf send labels and a logging.Debug message to Google cloud logging
func (labels Labels) Debugf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, labels.parse(), logging.Debug, service.LogTime())
}

// Warning send labels and a logging.Warning message to Google cloud logging
func (labels Labels) Warning(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, labels.parse(), logging.Warning, service.LogTime())
}

// Warningf send labels and a logging.Warning message to Google cloud logging
func (labels Labels) Warningf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, labels.parse(), logging.Warning, service.LogTime())
}

// Error send labels and a logging.Error message to Google cloud logging
func (labels Labels) Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, labels.parse(), logging.Error, service.LogTime())
}

// Errorf send labels and a logging.Error message to Google cloud logging
func (labels Labels) Errorf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, labels.parse(), logging.Error, service.LogTime())
}

// Fatal This is equivalent to calling Print followed by os.Exit(1)
// logged as logging.Alert
func (labels Labels) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, labels.parse(), logging.Alert, service.LogTime())
	BuiltInExit(1)
}

// Fatalf This is equivalent to calling Printf followed by os.Exit(1)
// logged as logging.Alert
func (labels Labels) Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, labels.parse(), logging.Alert, service.LogTime())
	BuiltInExit(1)
}

// Panic This is equivalent to calling Print followed by panic()
// logged as logging.Critical
func (labels Labels) Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, nil, nil, logging.Critical, service.LogTime())
	BuiltInPanic(message)
}

// Panicf This is equivalent to calling Printf followed by panic()
// logged as logging.Critical
func (labels Labels) Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, nil, nil, logging.Critical, service.LogTime())
	BuiltInPanic(message)
}

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
