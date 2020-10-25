package service

import (
	"github.com/alt4dev/protobuff/proto"
	"log"
	"testing"
)

var writeMock func(msg *proto.Log)

type remoteHelperMock struct{}

func (helper remoteHelperMock) WriteLog(msg *proto.Log, result *LogResult) {
	writeMock(msg)
	result.R = &proto.Result{
		Status:  proto.Result_UNAUTHORIZED,
		Message: "Mocked",
	}
	result.Err = nil
}

func (helper remoteHelperMock) WriteAudit(msg *proto.AuditLog, result *LogResult) {}

func (helper remoteHelperMock) QueryAudit(query proto.Query) (result *proto.QueryResult, err error) {
	return nil, nil
}

func TestOverrideDefaultLogOutputSync(t *testing.T) {
	syncLogger := log.New(SyncWriter, "alt4async ", 0)
	Alt4RemoteHelper = remoteHelperMock{}

	writeMock = func(msg *proto.Log) {
		if msg.Message != "alt4async This is a sync test i.e. will wait for result\n" {
			t.Error("Unexpected message logged")
			t.Error(msg.Message)
		}
	}
	syncLogger.Println("This is a sync test i.e. will wait for result")
}

func TestOverrideDefaultLogOutputAsync(t *testing.T) {
	// Be sure to wait for everything to finish
	defer WaitGroup().Wait()
	// The default logger writes to alt4
	asyncLogger := log.New(Writer, "alt4async ", 0)
	Alt4RemoteHelper = remoteHelperMock{}

	writeMock = func(msg *proto.Log) {
		if msg.Message != "alt4async This is an async test i.e. will write and forget\n" {
			t.Error("Unexpected message logged")
			t.Error(msg.Message)
		}
	}
	asyncLogger.Println("This is an async test i.e. will write and forget")
}
