package service
/*
Code used to write logs to the local machine usually notifying
- write failures to ALT4
- unauthenticated requests
- failed connections, e.t.c.
*/

import (
	"log"
)

var (
	emitError   *log.Logger
	emitWarning *log.Logger
	emit        *log.Logger
)

func init() {
	emitError = log.New(options.Writer, "[alt4](if seeing this please contact: critical@alt4.dev) ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	emitWarning = log.New(options.Writer, "[alt4] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	emit = log.New(options.Writer, "", 0)
}