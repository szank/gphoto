package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
import "C"
import "fmt"

//ContextLogCallback defineds a function used to log info associated to lobgphoto2 context
type ContextLogCallback func(string)

//LogCallback defines a generic libgphoto2 logging function
type LogCallback func(int, string, string)

// ContextInfoCallback is the function logging info logs from  libgphoto2 context.
//By default it logs everything to standard outout. You can assign your own method to this var
var ContextInfoCallback ContextLogCallback

// ContextErrorCallback is the function logging error logs from  libgphoto2 context.
//By default it logs everything to standard outout. You can assign your own method to this var
var ContextErrorCallback ContextLogCallback

// LoggerCallback is the libgphoto2 logging function. Currently there is no possibility to add multiple log function like it is possible in
// native C library implementation. Default implementation log everything to standard output with log level set to DEBUG
var LoggerCallback LogCallback

func defaultLoggerCallback(debugLevel int, domain, data string) {
	fmt.Println("LOGGING : domain " + domain + " data: " + data)
}

func defaultInfoCallback(data string) {
	fmt.Println("INFO: " + data)
}

func defaultErrorCallback(data string) {
	fmt.Println("ERROR : " + data)
}

//export wrapperInfoCallback
func wrapperInfoCallback(input *C.char) {
	if ContextInfoCallback != nil {
		ContextInfoCallback(C.GoString(input))
	}
}

//export wrapperErrorCallback
func wrapperErrorCallback(input *C.char) {
	if ContextErrorCallback != nil {
		ContextErrorCallback(C.GoString(input))
	}
}

//export wrapperLoggingCallback
func wrapperLoggingCallback(logLevel int, domain, data *C.char) {
	if LoggerCallback != nil {
		LoggerCallback(logLevel, C.GoString(domain), C.GoString(data))
	}
}

// make the general logging better. Make min log level settable somehow
func init() {
	ContextInfoCallback = defaultInfoCallback
	ContextErrorCallback = defaultErrorCallback
	LoggerCallback = defaultLoggerCallback
	C.gp_log_add_func(LogError, (*[0]byte)(C.loger_func), nil)
}
