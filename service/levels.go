package service


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
	CRITICAL: 4,
	FATAL:   5,
}

var levelString = map[uint8]string {
	LEVEL.INFO: "INFO",
	LEVEL.DEBUG: "DEBUG",
	LEVEL.WARNING: "WARNING",
	LEVEL.ERROR: "ERROR",
	LEVEL.CRITICAL: "PANIC",
	LEVEL.FATAL: "FATAL",
}