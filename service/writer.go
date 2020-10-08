package service


type alt4Writer struct {}

func (writer alt4Writer) Write(p []byte) (n int, err error) {
	message := string(p)
	_, _err := Log(3, false, message, nil, 1).Result()
	return len(p), _err
}

var Writer = alt4Writer{}
