// Package service holds the methods necessary to support implementation of loggers that write to Google cloud service
package service

import (
	"cloud.google.com/go/logging"
	"context"
	"io"
	"os"
	"sync"
	"time"
)

var options = struct {
	Mode   string
	Writer io.Writer
}{
	Mode:   "release",
	Writer: os.Stderr,
}

var client *logging.Client

func init() {
	SetMode(os.Getenv("ALT4_MODE"))

	// Don't create a client during testing or silent Modes
	if options.Mode == ModeTesting || options.Mode == ModeSilent {
		return
	}

	project := os.Getenv("PROJECT_ID")

	var err error
	client, err = logging.NewClient(context.Background(), project)
	if err != nil {
		panic(err)
	}

	client.OnError = func(err error) {
		emitError.Println(err)
	}
}

func setupOptions() {

}

const ModeRelease = "release"
const ModeDebug = "debug"
const ModeTesting = "testing"
const ModeSilent = "silent"

// SetMode Sets the behaviour of alt4 based on the following:
// `release` - Under this mode logs are written to Google cloud and not emitted to stdout
// `debug` - Under this mode logs are written to Google cloud and emitted to stdout
// `testing` - Under this mode logs are not written to Google cloud, just emitted to stdout
// `silent` - Under this mode logs are not written to Google cloud or emitted to stdout
// Mode can also be set using environment variable ALT4_MODE
// Default mode is `release`
func SetMode(mode string) {
	if mode == ModeRelease || mode == ModeDebug || mode == ModeTesting || mode == ModeSilent {
		options.Mode = mode
	}
}

// SetDebugOutput Is used to specify where alt4 emits additional output e.g. when facing network errors.
// Defaults os.Stderr
func SetDebugOutput(w io.Writer) {
	options.Writer = w
}

var timeLock sync.Mutex

func LogTime() time.Time {
	timeLock.Lock()
	defer timeLock.Unlock()
	return time.Now()
}
