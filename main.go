package gphoto

// #cgo LDFLAGS:  -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
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
	fmt.Println("Before capture image")
	path, err := camera.CaptureImage()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Printf("Captured file to %s/%s\n", path.Folder, path.Name)
	}

	var file *os.File
	file, err = os.Create("./" + path.Name)
	defer file.Close()

	if err != nil {
		fmt.Printf("Could not create file %s, error : %s\n", "./"+path.Name, err.Error())
	}

	///	if err = camera.DownloadImage(path, file, true); err != nil {
	//		fmt.Printf("Could not download file %s, error : %s\n", "./"+path.Name, err.Error())
	//	}

}
