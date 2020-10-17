// Package log provides methods for writing logs to Alt4
// Methods in this package write logs asynchronously unless otherwise specified.
// You can use the Package sync/log to write synchronously or call the function `Result` which will block if the operation is not done.
package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
	"os"
)

var LEVEL = service.LEVEL

// BuiltInPanic Internally this function just calls panic(). Override for testing(Panic, Panicf, Panicln)
var BuiltInPanic func(v interface{}) = func(v interface{}) {
	panic(v)
}

// BuiltInExit Internally this function just calls os.Exit. Override for testing(Fatal, Fatalf, Fatalln)
var BuiltInExit func(code int) = func(code int) {
	os.Exit(code)
}

// Group start a log group for the goroutine that calls this function.
// A group should be closed after. Use: `defer Group(...).Close()`
func Group(v ...interface{}) *GroupResult {
	title := fmt.Sprint(v...)
	return &GroupResult{
		logResult: service.Log(2, true, title, nil, LEVEL.DEBUG),
	}
}

// Print send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprint(a...)
func Print(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

// Printf send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprintf(a...)
func Printf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

// Println send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprintln(a...)
func Println(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

// Info send a log message to alt4. The log level is INFO. Log message will be formatted by fmt.Sprint(a...)
func Info(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO)
}

// Infof send a log message to alt4. The log level is INFO. Log message will be formatted by fmt.Sprintf(a...)
func Infof(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.INFO)
}

// Infoln send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprintln(a...)
func Infoln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO)
}

// Debug send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprint(a...)
func Debug(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

// Debugf send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprintf(a...)
func Debugf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

// Debugln send a log message to alt4. The log level is DEBUG. Log message will be formatted by fmt.Sprintln(a...)
func Debugln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

// Warning send a log message to alt4. The log level is WARNING. Log message will be formatted by fmt.Sprint(a...)
func Warning(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING)
}

// Warningf send a log message to alt4. The log level is WARNING. Log message will be formatted by fmt.Sprintf(a...)
func Warningf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING)
}

// Warningln send a log message to alt4. The log level is WARNING. Log message will be formatted by fmt.Sprintln(a...)
func Warningln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING)
}

// Error send a log message to alt4. The log level is ERROR. Log message will be formatted by fmt.Sprint(a...)
func Error(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR)
}

// Errorf send a log message to alt4. The log level is ERROR. Log message will be formatted by fmt.Sprintf(a...)
func Errorf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR)
}

// Errorln send a log message to alt4. The log level is ERROR. Log message will be formatted by fmt.Sprintln(a...)
func Errorln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR)
}

// Fatal This is equivalent to calling Print followed by os.Exit(1). The log level is CRITICAL.
func Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	BuiltInExit(1)
}

// Fatalf This is equivalent to calling Printf followed by os.Exit(1). The log level is CRITICAL.
func Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	BuiltInExit(1)
}

// Fatalln This is equivalent to calling Println followed by os.Exit(1). The log level is CRITICAL.
func Fatalln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	BuiltInExit(1)
}

// Panic This is equivalent to calling Print followed by panic(). The log level is CRITICAL.
func Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	BuiltInPanic(message)
}

// Panicf This is equivalent to calling Printf followed by panic(). The log level is CRITICAL.
func Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	BuiltInPanic(message)
}

// Panicln This is equivalent to calling Println followed by panic(). The log level is CRITICAL.
func Panicln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	BuiltInPanic(message)
}
