// Package service holds the methods necessary to support implementation of loggers that write to Google cloud service
package service

import (
	"cloud.google.com/go/logging"
	"context"
	"google.golang.org/genproto/googleapis/api/monitoredres"
	"io"
	"os"
	"sync"
	"time"
)

var options = struct {
	Mode     string
	Writer   io.Writer
	Resource *monitoredres.MonitoredResource
}{
	Mode:     "release",
	Writer:   os.Stderr,
	Resource: nil,
}

var client *logging.Client

func init() {
	SetMode(os.Getenv("ALT4_MODE"))
}

func setupClient() {
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

	// Set up the client after setting the mode
	setupClient()
}

// SetDebugOutput Is used to specify where alt4 emits additional output e.g. when facing network errors.
// Defaults os.Stderr
func SetDebugOutput(w io.Writer) {
	options.Writer = w
}

// SetMonitoredResource sets the monitored resource that is the source of log entries
// This accepts the type e.g. "gae_app", "cloud_run_revision" and labels which are simply a map of the resource labels
// The easiest way to get these values is to check a previous log entry printed by the same resource to CGP logging
func SetMonitoredResource(resourceType string, labels map[string]string) {
	options.Resource = &monitoredres.MonitoredResource{
		Type:   resourceType,
		Labels: labels,
	}
}

var timeLock sync.Mutex

func LogTime() time.Time {
	timeLock.Lock()
	defer timeLock.Unlock()
	return time.Now()
}
