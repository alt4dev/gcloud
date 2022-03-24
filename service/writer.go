package service

import (
	"cloud.google.com/go/logging"
)

type alt4Writer struct{}

func (writer alt4Writer) Write(p []byte) (n int, err error) {
	t := LogTime()
	message := string(p)

	Log(5, message, nil, nil, logging.Default, t)
	return len(p), nil
}

// Writer can be used to override a normal/default go logger to write it's output to alt4
// This writer still respects log grouping if you're grouping your request logs
// Example log.SetOutput(Writer) for the default log package
var Writer = alt4Writer{}
