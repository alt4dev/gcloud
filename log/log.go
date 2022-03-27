// Package log provides methods for writing logs to Google cloud logging
// The biggest edge this library brings is grouping logs, it's advisable to take advantage of this behaviour
package log

import (
	"cloud.google.com/go/logging"
	"fmt"
	"github.com/alt4dev/gcloud/service"
	"os"
)

// BuiltInPanic Internally this function just calls panic(). Override for testing(Panic, Panicf)
var BuiltInPanic func(v interface{}) = func(v interface{}) {
	panic(v)
}

// BuiltInExit Internally this function just calls os.Exit. Override for testing(Fatal, Fatalf)
var BuiltInExit func(code int) = func(code int) {
	os.Exit(code)
}

// Print send a log message to Google cloud logging
func Print(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Default, service.LogTime())
}

// Printf send a log message to Google cloud logging
func Printf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Default, service.LogTime())
}

// Println send a log message to Google cloud logging
func Println(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, message, logging.Default, service.LogTime())
}

// Info send a logging.Info message to Google cloud logging
func Info(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Info, service.LogTime())
}

// Infof send a logging.Info message to Google cloud logging
func Infof(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Info, service.LogTime())
}

// Debug send a logging.Debug message to Google cloud logging
func Debug(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Debug, service.LogTime())
}

// Debugf send a logging.Debug message to Google cloud logging
func Debugf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Debug, service.LogTime())
}

// Warning send a logging.Warning message to Google cloud logging
func Warning(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Warning, service.LogTime())
}

// Warningf send a logging.Warning message to Google cloud logging
func Warningf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Warning, service.LogTime())
}

// Error send a logging.Error message to Google cloud logging
func Error(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Error, service.LogTime())
}

// Errorf send a logging.Error message to Google cloud logging
func Errorf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Error, service.LogTime())
}

// Fatal This is equivalent to calling Print followed by os.Exit(1)
// logged as logging.Alert
func Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Alert, service.LogTime())
	service.CloseGroup(0, nil)
	BuiltInExit(1)
}

// Fatalf This is equivalent to calling Printf followed by os.Exit(1)
// logged as logging.Alert
func Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Alert, service.LogTime())
	service.CloseGroup(0, nil)
	BuiltInExit(1)
}

// Panic This is equivalent to calling Print followed by panic()
// logged as logging.Critical
func Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, message, logging.Critical, service.LogTime())
	BuiltInPanic(message)
}

// Panicf This is equivalent to calling Printf followed by panic()
// logged as logging.Critical
func Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, message, logging.Critical, service.LogTime())
	BuiltInPanic(message)
}
