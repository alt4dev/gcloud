package service

import (
	"encoding/json"
	"github.com/alt4dev/protobuff/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"io"
	"io/ioutil"
	"os"
)

var address = "rpc.alt4.dev:443"
var connection *grpc.ClientConn
var client *proto.LoggingClient

type Alt4Options struct {
	AuthToken string
	Mode      string
	Sink      string
	Writer    io.Writer
}

var options = Alt4Options{
	AuthToken: "",
	Mode:      "release",
	Sink:      "default",
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
			EmitError.Println("Error creating a connection to alt4. Error: ", err)
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
			EmitError.Print("Error opening file provided in ALT4_CONFIG. Error: ", err)
		}else {
			content := struct {
				Token     string `json:"token" binding:"required"`
				Mode      string `json:"mode"`
				Sink		string `json:"sink"`
			}{}
			err = json.Unmarshal(jsonContent, &content)
			if err != nil {
				EmitError.Println("Error decoding ALT4_CONFIG. Error: ", err)
			}else {
				SetAuthToken(content.Token)
				SetMode(content.Mode)
				SetSink(content.Sink)
			}
		}
	}
	SetAuthToken(os.Getenv("ALT4_AUTH_TOKEN"))
	SetMode(os.Getenv("ALT4_MODE"))
	SetSink(os.Getenv("ALT4_SINK"))
}

func SetAuthToken(token string) {
	if token != "" {
		options.AuthToken = token
	}
}

func SetMode(mode string) {
	if mode == "release" || mode == "debug" {
		options.Mode = mode
	}
}

func SetSink(sink string) {
	if sink != "" {
		options.Sink = sink
	}
}

func SetDebugOutput(w io.Writer){
	options.Writer = w
}
