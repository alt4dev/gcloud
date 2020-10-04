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
	EmitError   *log.Logger
	EmitWarning *log.Logger
	Emit        *log.Logger
)

func init() {
	EmitError = log.New(options.EmmitFile, "[alt4](if seeing this please contact: critical@alt4.dev) ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	EmitWarning = log.New(options.EmmitFile, "[alt4] WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	Emit = log.New(options.EmmitFile, "[alt4emmit] ", log.Ldate|log.Ltime|log.Lshortfile)
}