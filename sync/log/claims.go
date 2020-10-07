package log

import (
	"fmt"
	alt4log "github.com/alt4dev/go/log"
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
	"log"
)
 
type Claims map[string]interface{}

func (claims Claims) parse() []*proto.Claim {
	return alt4log.ParseClaims(alt4log.Claims(claims))
}

func (claims Claims) Print(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, claims.parse(), LEVEL.DEBUG).Result()
}

func (claims Claims) Printf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, claims.parse(), LEVEL.DEBUG).Result()
}

func (claims Claims) Println(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, claims.parse(), LEVEL.DEBUG).Result()
}

func (claims Claims) Info(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO).Result()
}

func (claims Claims) Infof(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.INFO).Result()
}

func (claims Claims) Infoln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.INFO).Result()
}

func (claims Claims) Debug(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func (claims Claims) Debugf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func (claims Claims) Debugln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.DEBUG).Result()
}

func (claims Claims) Warning(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING).Result()
}

func (claims Claims) Warningf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING).Result()
}

func (claims Claims) Warningln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.WARNING).Result()
}

func (claims Claims) Error(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprint(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR).Result()
}

func (claims Claims) Errorf(format string, v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintf(format, v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR).Result()
}

func (claims Claims) Errorln(v ...interface{}) (*proto.Result, error) {
	message := fmt.Sprintln(v...)
	return service.Log(2, false, message, nil, LEVEL.ERROR).Result()
}

func (claims Claims) Fatal(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.ERROR).Result()
	log.Fatal(v...)
}

func (claims Claims) Fatalf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.ERROR).Result()
	log.Fatalf(format, v...)
}

func (claims Claims) Fatalln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.ERROR).Result()
	log.Fatalln(v...)
}

func (claims Claims) Panic(v ...interface{}) {
	message := fmt.Sprint(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	log.Panic(v...)
}

func (claims Claims) Panicf(format string, v ...interface{}) {
	message := fmt.Sprintf(format, v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	log.Panicf(format, v...)
}

func (claims Claims) Panicln(v ...interface{}) {
	message := fmt.Sprintln(v...)
	service.Log(2, false, message, nil, LEVEL.CRITICAL).Result()
	log.Panicln(v...)
}
