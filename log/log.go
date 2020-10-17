package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
	"os"
)

var LEVEL = service.LEVEL

// Override for testing
var builtInPanic func(v interface{}) = func(v interface{}) {
	panic(v)
}

// Override for testing
var builtInExit func(code int) = func(code int) {
	os.Exit(code)
}

func Print(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

func Printf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

func Println(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

func Info(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO)
}

func Infof(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.INFO)
}

func Infoln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO)
}

func Debug(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

func Debugf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

func Debugln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG)
}

func Warning(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING)
}

func Warningf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING)
}

func Warningln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING)
}

func Error(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR)
}

func Errorf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR)
}

func Errorln(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR)
}

func Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	builtInExit(1)
}

func Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	builtInExit(1)
}

func Fatalln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	builtInExit(1)
}

func Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	builtInPanic(message)
}

func Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	builtInPanic(message)
}

func Panicln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	builtInPanic(message)
}
