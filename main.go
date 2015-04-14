package main

// #cgo LDFLAGS: -L. -lgphotocallbacks -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
import "C"
import (
	"fmt"
	"os"
)

func main() {
	context, err := GetNewGPhotoContext()
	if err != nil {
		fmt.Printf("Error creating context !\n\n")
		return
	}
	//TOODO: add finalizer
	defer context.Free()

	camera, err := GetNewGPhotoCamera(context)
	if err != nil {
		fmt.Printf("Error initializing camera : %v", err.Error())
		return
	}

	camera.GetWidgetTree()
	camera.PrintWidgetTree(os.Stdout)
}

func init() {
	C.gp_log_add_func(LogDebug, (*[0]byte)(C.loger_func), nil)
}
