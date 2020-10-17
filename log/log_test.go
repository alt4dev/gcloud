package log

import (
	"fmt"
	"github.com/alt4dev/go/service"
	"github.com/alt4dev/protobuff/proto"
	"runtime"
	"testing"
)

var writerMock func(msg *proto.Message)

type RemoteWriterMock struct{}

func (w RemoteWriterMock) Write(msg *proto.Message, result *service.LogResult) {
	writerMock(msg)
}

func whereAmI() int {
	_, _, l, _ := runtime.Caller(1)
	return l
}

var testFile string
var testMessage string
var testLine int

func setUp(t *testing.T, level uint8, isGroup bool, claims []*proto.Claim) {
	service.Alt4RemoteWriter = RemoteWriterMock{}

	_, testFile, _, _ = runtime.Caller(1)
	writerMock = func(msg *proto.Message) {
		if msg.Message != testMessage {
			t.Error("Unexpected testMessage found")
			t.Log(msg.Message)
			t.Log(testMessage)
		}
		if msg.LineNo != uint32(testLine) {
			t.Errorf("Unexpected testLine  number found. %d != %d\n", msg.LineNo, testLine)
		}
		if msg.FileName != testFile {
			t.Errorf("Unexpected testFile found. '%s' != '%s'", msg.FileName, testFile)
		}
		if msg.Level != uint32(level) {
			t.Errorf("Unexpected level found. %d != %d\n", msg.Level, level)
		}
		if msg.ThreadInit != isGroup {
			t.Errorf("Expected thread init to be %v but found %v", isGroup, msg.ThreadInit)
		}
		if fmt.Sprint(claims) != fmt.Sprint(msg.Claims) {
			t.Errorf("Claims don't match")
			t.Log("Expected: ", claims)
			t.Log("Found: ", msg.Claims)
		}
	}
}


func TestGroup(t *testing.T) {
	setUp(t, LEVEL.DEBUG, true, nil)
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	defer Group("A test print testMessage", "nothing", 10).Close()
}

func TestPrint(t *testing.T) {
	setUp(t, LEVEL.DEBUG, false, nil)
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = Print("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = Println("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = Printf("A test formatting, here's a number %d", 145).Result()
}

func TestInfo(t *testing.T) {
	setUp(t, LEVEL.INFO, false, nil)
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = Info("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = Infoln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = Infof("A test formatting, here's a number %d", 145).Result()
}

func TestDebug(t *testing.T) {
	setUp(t, LEVEL.DEBUG, false, nil)
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = Debug("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = Debugln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = Debugf("A test formatting, here's a number %d", 145).Result()
}

func TestWarning(t *testing.T) {
	setUp(t, LEVEL.WARNING, false, nil)
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = Warning("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = Warningln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = Warningf("A test formatting, here's a number %d", 145).Result()
}

func TestError(t *testing.T) {
	setUp(t, LEVEL.ERROR, false, nil)
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = Error("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = Errorln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = Errorf("A test formatting, here's a number %d", 145).Result()
}

func TestFatal(t *testing.T) {
	setUp(t, LEVEL.CRITICAL, false, nil)
	// Mock exit
	BuiltInExit = func(code int) {
		if code != 1 {
			t.Errorf("Fatal should call os.Exit with exit code 1. %d!=1", code)
		}
	}
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	Fatal("A test print testMessage", "nothing", 10)
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	Fatalln("A test new testLine", "something", 16)
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	Fatalf("A test formatting, here's a number %d", 145)
}

func TestPanic(t *testing.T) {
	setUp(t, LEVEL.CRITICAL, false, nil)
	// Mock exit
	BuiltInPanic = func(v interface{}) {
		if v.(string) != testMessage {
			t.Errorf("Unexpected panic testMessage. '%s' != '%s'", v.(string), testMessage)
		}
	}
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	Panic("A test print testMessage", "nothing", 10)
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	Panicln("A test new testLine", "something", 16)
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	Panicf("A test formatting, here's a number %d", 145)
}

