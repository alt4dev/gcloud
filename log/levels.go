package log


var LEVEL = struct {
	INFO uint8
	DEBUG uint8
	WARNING uint8
	ERROR uint8
	FATAL uint8
	CRITICAL uint8
}{
	INFO:    0,
	DEBUG:   1,
	WARNING: 2,
	ERROR:   3,
	FATAL:   4,
	CRITICAL: 5,
}
