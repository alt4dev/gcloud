package service

import (
	"github.com/alt4dev/protobuff/proto"
	"log"
	"testing"
)

var writeMock func(msg *proto.Message)
type RemoteWriterMock struct {}

func (w RemoteWriterMock) Write(msg *proto.Message, result *LogResult) {
	writeMock(msg)
	result.R = &proto.Result{
		Acknowledged: true,
		Message:      "Mocked",
	}
	result.Err = nil
}

func TestOverrideDefaultLogOutputSync(t *testing.T) {
	syncLogger := log.New(SyncWriter, "alt4async ", 0)
	Alt4RemoteWriter = RemoteWriterMock{}

	writeMock = func(msg *proto.Message) {
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
	Alt4RemoteWriter = RemoteWriterMock{}

	writeMock = func(msg *proto.Message) {
		if msg.Message != "alt4async This is an async test i.e. will write and forget\n" {
			t.Error("Unexpected message logged")
			t.Error(msg.Message)
		}
	}
	asyncLogger.Println("This is an async test i.e. will write and forget")


}