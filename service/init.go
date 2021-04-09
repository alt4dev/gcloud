// Package service holds the methods necessary to support implementation of loggers that write to alt4.
package service

import (
	"context"
	"encoding/json"
	"github.com/alt4dev/protobuff/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"io"
	"io/ioutil"
	"os"
	"sync"
	"time"
)

var address = "rpc.alt4.dev:443"
var connection *grpc.ClientConn
var client *proto.LoggingClient

var options = struct {
	AuthContext context.Context
	Mode      string
	Source      string
	Writer    io.Writer
}{
	AuthContext: nil,
	Mode:      "release",
	Source:      "default",
	Writer:    os.Stderr,
}

func init() {
	// Initialize client on connection
	getClient()
	// Setup options from env
	setupOptions()
}

// getClient connects to alt4 and creates a new logging client or reuses an existing connection
func getClient() *proto.LoggingClient {
	if connection == nil || client == nil {
		var err error
		transCert := grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(nil, "rpc.alt4.dev"))
		connection, err = grpc.Dial(address, transCert)
		if err != nil {
			emitError.Println("Error creating a connection to alt4. Error: ", err)
			connection = nil
			client = nil
			return nil
		}
		c := proto.NewLoggingClient(connection)
		client = &c
	}
	return client
}

func setupOptions() {
	optionsFile := os.Getenv("ALT4_CONFIG")
	if optionsFile != "" {
		jsonContent, err := ioutil.ReadFile(optionsFile)
		if err != nil {
			emitError.Printf("Error opening file `%s` provided in ALT4_CONFIG. Error: %s\n", optionsFile, err)
		}else {
			content := struct {
				Token     string `json:"token" binding:"required"`
				Mode      string `json:"mode"`
				Source		string `json:"source"`
			}{}
			err = json.Unmarshal(jsonContent, &content)
			if err != nil {
				emitError.Println("Error decoding ALT4_CONFIG. Error: ", err)
			}else {
				SetAuthToken(content.Token)
				SetMode(content.Mode)
				SetSource(content.Source)
			}
		}
	}
	SetAuthToken(os.Getenv("ALT4_AUTH_TOKEN"))
	SetMode(os.Getenv("ALT4_MODE"))
	SetSource(os.Getenv("ALT4_SOURCE"))
}

// SetAuthToken Used to set the auth token for writing to alt4.
// This setting can be done via config file ALT4_CONFIG or setting environment variable ALT4_AUTH_TOKEN
func SetAuthToken(token string) {
	md := metadata.Pairs("AuthToken", token)
	if token != "" {
		options.AuthContext = metadata.NewOutgoingContext(context.Background(), md)
	}
}

const ModeRelease = "release"
const ModeDebug = "debug"
const ModeTesting = "testing"
const ModeSilent = "silent"

// SetMode Sets the behaviour of alt4 based on the following:
// `release` - Under this mode logs are written to alt4 and not emitted to stdout
// `debug` - Under this mode logs are written to alt4 and emitted to stdout
// `testing` - Under this mode logs are not written to alt4, just emitted to stdout
// `silent` - Under this mode logs are not written to alt4 or emitted to stdout
// `json`(coming soon) - Under this mode all logs are written to a JSON file which you can later upload to alt4
// Mode can also be set via a config file ALT4_CONFIG or setting environment variable ALT4_MODE
// Default mode is `release`
func SetMode(mode string) {
	if mode == ModeRelease || mode == ModeDebug || mode == ModeTesting || mode == ModeSilent {
		options.Mode = mode
	}
}

// SetSource Sets an id to identify where your logs are coming from
// By setting the source you distinguish logs from different sources e.g. Languages, services, servers e.t.c.
func SetSource(source string) {
	if source != "" {
		options.Source = source
	}
}

// SetDebugOutput Is used to specify where alt4 emits additional output e.g. when facing network errors.
// Defaults os.Stderr
func SetDebugOutput(w io.Writer){
	options.Writer = w
}

var timeLock sync.Mutex

func LogTime() time.Time {
	timeLock.Lock()
	defer timeLock.Unlock()
	return time.Now()
}
