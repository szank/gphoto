package main

// #cgo LDFLAGS: -L. -lgphotocallbacks -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
// #include <stdlib.h>
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

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

type cameraFilePathInternal struct {
	Name   [128]uint8
	Folder [1024]uint8
}

//CameraFilePath is a path to a file on the camera file system
type CameraFilePath struct {
	Name   string
	Folder string
}

func newCameraFilePathFromInternalImpl(input *cameraFilePathInternal) *CameraFilePath {
	return &CameraFilePath{
		Name:   string(input.Name[:bytes.IndexByte(input.Name[:], 0)]),
		Folder: string(input.Folder[:bytes.IndexByte(input.Folder[:], 0)]),
	}
}

//GetNewGPhotoCamera returns a new camera instance
func GetNewGPhotoCamera(context *Context) (*Camera, error) {
	var gpCamera *C.Camera
	C.gp_camera_new((**C.Camera)(unsafe.Pointer(&gpCamera)))

	if gpCamera == nil {
		return nil, fmt.Errorf("Cannot initialize camera pointer")
	}

	initCode := C.gp_camera_init(gpCamera, context.gpContext)
	if initCode != gpOk {
		C.gp_camera_exit(gpCamera, context.gpContext)
		C.gp_camera_unref(gpCamera)
		return nil, fmt.Errorf("Error number : %d", initCode)
	}
	return &Camera{
		gpCamera:  gpCamera,
		gpContext: context.gpContext,
	}, nil
}

//GetWidgetTree returns a widget tree for selected camera
func (camera *Camera) GetWidgetTree() error {
	var rootWidget *C.CameraWidget

	if retval := C.gp_camera_get_config(camera.gpCamera, (**C.CameraWidget)(unsafe.Pointer(&rootWidget)), camera.gpContext); retval != gpOk {
		return fmt.Errorf("cannot initialize camera settings tree error code :%v", retval)
	}
	defer C.gp_widget_free(rootWidget)
	camera.getWidgetInfo(rootWidget, &camera.CameraSettings)
	return nil
}

func (camera *Camera) getWidgetInfo(input *C.CameraWidget, output *CameraWidget) {
	var child *C.CameraWidget
	var data *C.char
	var widgetType C.CameraWidgetType

	//ignore return value (int) for now
	C.gp_widget_get_info(input, (**C.char)(unsafe.Pointer(&data)))
	output.Info = C.GoString(data)

	C.gp_widget_get_label(input, (**C.char)(unsafe.Pointer(&data)))
	output.Label = C.GoString(data)

	C.gp_widget_get_name(input, (**C.char)(unsafe.Pointer(&data)))
	output.Label = C.GoString(data)

	C.gp_widget_get_type(input, (*C.CameraWidgetType)(unsafe.Pointer(&widgetType)))
	output.Type = int(widgetType)

	childrenCount := int(C.gp_widget_count_children(input))
	output.Children = make([]CameraWidget, childrenCount)
	for n := 0; n < childrenCount; n++ {
		C.gp_widget_get_child(input, C.int(n), (**C.CameraWidget)(unsafe.Pointer(&child)))
		camera.getWidgetInfo(child, &output.Children[n])
	}
}

//PrintWidgetTree prints widget hierarchy to the output buffer
func (camera *Camera) PrintWidgetTree(file io.Writer) {

	widget := &camera.CameraSettings
	fmt.Printf("Widget Info[%v], Label[%v] , Name[%v], type [%v]\n", widget.Info, widget.Label, widget.Name, widget.Type)
	for _, child := range widget.Children {
		fmt.Printf("    Widget Info[%v], Label[%v] , Name[%v], type [%v]\n", child.Info, child.Label, child.Name, child.Type)
	}

}

//CaptureImage captures image with current setings into camera's internal storage
func (camera *Camera) CaptureImage() (*CameraFilePath, error) {
	photoPath := cameraFilePathInternal{}

	if retval := C.gp_camera_capture(camera.gpCamera, 0, (*C.CameraFilePath)(unsafe.Pointer(&photoPath)), camera.gpContext); retval != gpOk {
		return nil, fmt.Errorf("Cannot capture photo, error code :%v", retval)
	}
	return newCameraFilePathFromInternalImpl(&photoPath), nil
}

//DownloadImage saves image pointed by path to the provided buffer. If leave on camera is set to false,the file will be deleted from the camera internal storage
func (camera *Camera) DownloadImage(path *CameraFilePath, buffer io.Writer, leaveOnCamera bool) error {
	var gpFile *C.CameraFile
	C.gp_file_new((**C.CameraFile)(unsafe.Pointer(&gpFile)))

	if gpFile == nil {
		return fmt.Errorf("Cannot initialize camera file")
	}
	defer C.gp_file_free(gpFile)

	fileDir := C.CString(path.Folder)
	defer C.free(unsafe.Pointer(fileDir))

	fileName := C.CString(path.Name)
	defer C.free(unsafe.Pointer(fileName))

	if retval := C.gp_camera_file_get(camera.gpCamera, fileDir, fileName, FileTypeNormal, gpFile, camera.gpContext); retval != gpOk {
		return fmt.Errorf("Cannot download photo file, error code :%v", retval)
	}

	var fileData *C.char
	var fileLen C.ulong
	C.gp_file_get_data_and_size(gpFile, (**C.char)(unsafe.Pointer(&fileData)), &fileLen)

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(fileData)),
		Len:  int(fileLen),
		Cap:  int(fileLen),
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	_, err := buffer.Write(goSlice)
	if err != nil && leaveOnCamera == false {
		C.gp_camera_file_delete(camera.gpCamera, fileDir, fileName, camera.gpContext)
	}
	return err
}
