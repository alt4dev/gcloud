package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
)

func Print(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(false, message, nil, LEVEL.DEBUG)
}

func Printf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(false, message, nil, LEVEL.DEBUG)
}

func Println(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(false, message, nil, LEVEL.DEBUG)
}

