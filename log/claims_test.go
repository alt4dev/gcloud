package log

import (
	"fmt"
	"testing"
)

var testClaims = Claims{
	"name": "Tester",
	"age": 25,
	"notStupid": false,
}

func TestClaims_Group(t *testing.T) {
	setUp(t, LEVEL.DEBUG, true, testClaims.parse())
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	defer testClaims.Group("A test print testMessage", "nothing", 10).Close()
}

func TestClaims_Print(t *testing.T) {
	setUp(t, LEVEL.DEBUG, false, testClaims.parse())
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = testClaims.Print("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = testClaims.Println("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = testClaims.Printf("A test formatting, here's a number %d", 145).Result()
}

func TestClaims_Info(t *testing.T) {
	setUp(t, LEVEL.INFO, false, testClaims.parse())
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = testClaims.Info("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = testClaims.Infoln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = testClaims.Infof("A test formatting, here's a number %d", 145).Result()
}

func TestClaims_Debug(t *testing.T) {
	setUp(t, LEVEL.DEBUG, false, testClaims.parse())
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = testClaims.Debug("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = testClaims.Debugln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = testClaims.Debugf("A test formatting, here's a number %d", 145).Result()
}

func TestClaims_Warning(t *testing.T) {
	setUp(t, LEVEL.WARNING, false, testClaims.parse())
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = testClaims.Warning("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = testClaims.Warningln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = testClaims.Warningf("A test formatting, here's a number %d", 145).Result()
}

func TestClaims_Error(t *testing.T) {
	setUp(t, LEVEL.ERROR, false, testClaims.parse())
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	_, _ = testClaims.Error("A test print testMessage", "nothing", 10).Result()
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	_, _ = testClaims.Errorln("A test new testLine", "something", 16).Result()
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	_, _ = testClaims.Errorf("A test formatting, here's a number %d", 145).Result()
}

func TestClaims_Fatal(t *testing.T) {
	setUp(t, LEVEL.CRITICAL, false, testClaims.parse())
	// Mock exit
	BuiltInExit = func(code int) {
		if code != 1 {
			t.Errorf("Fatal should call os.Exit with exit code 1. %d!=1", code)
		}
	}
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	testClaims.Fatal("A test print testMessage", "nothing", 10)
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	testClaims.Fatalln("A test new testLine", "something", 16)
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	testClaims.Fatalf("A test formatting, here's a number %d", 145)
}

func TestClaims_Panic(t *testing.T) {
	setUp(t, LEVEL.CRITICAL, false, testClaims.parse())
	// Mock exit
	BuiltInPanic = func(v interface{}) {
		if v.(string) != testMessage {
			t.Errorf("Unexpected panic testMessage. '%s' != '%s'", v.(string), testMessage)
		}
	}
	testMessage = fmt.Sprint("A test print testMessage", "nothing", 10)
	testLine = whereAmI() + 1
	testClaims.Panic("A test print testMessage", "nothing", 10)
	testMessage = fmt.Sprintln("A test new testLine", "something", 16)
	testLine = whereAmI() + 1
	testClaims.Panicln("A test new testLine", "something", 16)
	testMessage = fmt.Sprintf("A test formatting, here's a number %d", 145)
	testLine = whereAmI() + 1
	testClaims.Panicf("A test formatting, here's a number %d", 145)
}

