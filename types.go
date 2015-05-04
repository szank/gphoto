package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
import "C"

//Camera struct represents a camera connected to the computer
type Camera struct {
	gpCamera       *C.Camera
	gpContext      *C.GPContext
	CameraSettings CameraWidget
}

//CameraWidget is a representation of one of the cameras's setting or control knob
type CameraWidget struct {
	Label    string
	Name     string
	Info     string
	Type     int
	Children []CameraWidget
}

//Context represents a  context in which all other calls are executed
type Context struct {
	gpContext *C.GPContext
}

//CameraFilePath is a path to a file or dir on the camera file system
type CameraFilePath struct {
	Name     string
	Folder   string
	Isdir    bool
	Children []CameraFilePath

	camera *Camera
}

type CameraStorageInfo struct {
	Description string
	Capacity    uint64
	Free        uint64
	FreeImages  uint64
	Children    []CameraFilePath

	basedir string
	camera  *Camera
}
