package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
	"time"
)

type Claims map[string]interface{}

func (claims Claims) parse() []*proto.Claim {
	protoClaims := make([]*proto.Claim, 0)
	for key, i := range claims {
		var claimValue string
		var claimType uint8
		switch i.(type) {
		case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64:
			claimType = 1
			claimValue = fmt.Sprint(i)
		case float32, float64:
			claimType = 2
			claimValue = fmt.Sprint(i)
		case bool:
			claimType = 3
			claimValue = fmt.Sprint(i.(bool))
		case string:
			claimType = 4
			claimValue = i.(string)
		case time.Time:
			claimType = 5
			claimValue = fmt.Sprint(i.(time.Time).UnixNano())
		default:
			claimType = 0
			claimValue = fmt.Sprint(i)
		}
		protoClaims = append(protoClaims, &proto.Claim{
			Name:     key,
			DataType: uint32(claimType),
			Value:    claimValue,
		})
	}
	return protoClaims
}

func (claims Claims) Print(v ...interface{}) *service.LogResult {
	message := fmt.Sprint(v...)
	return service.Log(false, message, claims.parse(), LEVEL.DEBUG)
}

func (claims Claims) Printf(format string, v ...interface{}) *service.LogResult {
	message := fmt.Sprintf(format, v...)
	return service.Log(false, message, claims.parse(), LEVEL.DEBUG)
}

func (claims Claims) Println(v ...interface{}) *service.LogResult {
	message := fmt.Sprintln(v...)
	return service.Log(false, message, claims.parse(), LEVEL.DEBUG)
}
