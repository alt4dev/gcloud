package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
	"log"
)

var LEVEL = service.LEVEL

func Print(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func Printf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func Println(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func Info(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO).Result()
}

func Infof(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.INFO).Result()
}

func Infoln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO).Result()
}

func Debug(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func Debugf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func Debugln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func Warning(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING).Result()
}

func Warningf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING).Result()
}

func Warningln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING).Result()
}

func Error(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR).Result()
}

func Errorf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR).Result()
}

func Errorln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR).Result()
}

func Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.ERROR).Result()
	log.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.ERROR).Result()
	log.Fatalf(format, v...)
}

func Fatalln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.ERROR).Result()
	log.Fatalln(v...)
}

func Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	log.Panic(v...)
}

func Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	log.Panicf(format, v...)
}

func Panicln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	log.Panicln(v...)
}
