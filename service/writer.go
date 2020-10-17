package service

type alt4Writer struct{}

func (writer alt4Writer) Write(p []byte) (n int, err error) {
	message := string(p)
	Log(4, false, message, nil, 1)
	return len(p), nil
}

type alt4SyncWriter struct{}

func (writer alt4SyncWriter) Write(p []byte) (n int, err error) {
	message := string(p)
	_, _err := Log(4, false, message, nil, 1).Result()
	return len(p), _err
}

// Writer can be used to override a normal/default go logger to write it's output to alt4
// This method writes logs asynchronously. Opening and Closing a group at the end of your routines ensure waits for all writes to finish.
// Example log.SetOutput(Writer)
var Writer = alt4Writer{}

// SyncWriter can be used to override a normal/default go logger to write it's output to alt4
// This method writes logs synchronously.
// Unless you really need to we instead advice on Opening and Closing a group at the end of your routines which will wait for all writes started within the routine to complete
// Example log.SetOutput(SyncWriter)
var SyncWriter = alt4SyncWriter{}
