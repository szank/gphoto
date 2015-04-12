package main

import "C"
import "fmt"

//Callbacks part :
type ContextLogCallback func(string)
type LogCallback func(int, string, string)

var ContextInfoCallback ContextLogCallback
var ContextErrorCallback ContextLogCallback
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

func addLoggingFunc(LogCallback value) int

func init() {
	ContextInfoCallback = defaultInfoCallback
	ContextErrorCallback = defaultErrorCallback
	LoggerCallback = defaultLoggerCallback
}
