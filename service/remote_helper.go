package service

import (
	"errors"
	"github.com/alt4dev/protobuff/proto"
	"sync"
)

// LogResult Object returned when you create a log entry.
type LogResult struct {
	R   *proto.Result
	wg *sync.WaitGroup
	Err error
}

// Result Returns actual Result from alt4. This will block and wait for the Result if not done
func (result *LogResult) Result() (*proto.Result, error) {
	result.wg.Wait()
	return result.R, result.Err
}

// RemoteWriter an interface for functions called when writing to alt4.
// You can implement this function to mock writes to alt4 for better testing of your system.
type RemoteHelper interface {
	// WriteLog function will be called with the Message to be sent to alt4 and an empty LogResult to fill once done
	WriteLog(msg *proto.Log, result *LogResult)
	// WriteAudit function will be called with an audit message to be sent to alt4 and an empty LogResult to fill once done.
	WriteAudit(msg *proto.AuditLog, result *LogResult)
	// QueryAudit function will be called when you query audit logs. This function is synchronous.
	QueryAudit(query proto.Query) (result *proto.QueryResult, err error)
}

type DefaultHelper struct{}

func (helper DefaultHelper) WriteLog(msg *proto.Log, result *LogResult) {
	if getClient() == nil {
		result.Err = errors.New("error connecting to remote server")
		emitLog(msg, result.Err)
		return
	}
	if options.Mode != "testing" {
		result.R, result.Err = (*client).WriteLog(options.AuthContext, msg)
	}
	shouldEmit := options.Mode == "debug" || options.Mode == "testing"
	if (result.R != nil && result.R.Status != proto.Result_ACKNOWLEDGED) || shouldEmit || result.Err != nil {
		if result.R != nil && result.R.Status != proto.Result_ACKNOWLEDGED {
			result.Err = errors.New(result.R.Message)
		}
		emitLog(msg, result.Err)
	}
}

func (helper DefaultHelper) WriteAudit(msg *proto.AuditLog, result *LogResult){
	// TODO: Implement audit
}

func (helper DefaultHelper) QueryAudit(query proto.Query) (result *proto.QueryResult, err error) {
	// TODO: Implement query audit
	return nil, nil
}

// Alt4RemoteWriter For testing purposes, implement your own RemoteHelper and equate it to this variable
var Alt4RemoteHelper RemoteHelper = DefaultHelper{}
