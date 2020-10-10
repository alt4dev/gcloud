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

var Writer = alt4Writer{}
var SyncWriter = alt4SyncWriter{}
