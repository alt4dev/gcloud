package service

import (
	"context"
	"errors"
	"github.com/alt4dev/protobuff/proto"
)

// RemoteWriter an interface for functions called when writing to alt4.
// You can implement this function to mock writes to alt4 for better testing of your system.
type RemoteHelper interface {
	// WriteLog function will be called with the Message to be sent to alt4 and an empty R to fill once done
	WriteLog(msg *proto.Message, result *LogResult)

	WriteAudit(msg *proto.AuditMessage, result *LogResult)
}

type writer struct{}

func (w writer) WriteLog(msg *proto.Message, result *LogResult) {
	if getClient() == nil {
		result.Err = errors.New("error connecting to remote server")
		emitLog(msg, result.Err)
		return
	}
	if options.Mode != "testing" {
		result.R, result.Err = (*client).Log(context.Background(), msg)
	}
	shouldEmit := options.Mode == "debug" || options.Mode == "testing"
	if (result.R != nil && result.R.Status != proto.Result_ACKNOWLEDGED) || shouldEmit || result.Err != nil {
		if result.R != nil && result.R.Status != proto.Result_ACKNOWLEDGED {
			result.Err = errors.New(result.R.Message)
		}
		emitLog(msg, result.Err)
	}
}

func (w writer)