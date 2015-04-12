package main

// #cgo LDFLAGS: -L. -lgphotocallbacks -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
import "C"
import (
	"fmt"
	"io"
	"unsafe"
)

type Camera struct {
	gpCamera       *C.Camera
	gpContext      *C.GPContext
	CameraSettings CameraWidget
}

type CameraWidget struct {
	Label    string
	Name     string
	Info     string
	Type     int
	Children []CameraWidget
}

func GetNewGPhotoCamera() (*Camera, error) {
	var gpCamera *C.Camera
	C.gp_camera_new((**C.Camera)(unsafe.Pointer(&gpCamera)))

	if gpCamera == nil {
		return nil, fmt.Errorf("Cannot initialize camera pointer")
	}
	return &Camera{
		gpCamera: gpCamera,
	}, nil

}

func (camera *Camera) Init(context *Context) error {
	camera.gpContext = context.gpContext

	retval := C.gp_camera_init(camera.gpCamera, camera.gpContext)
	if retval != GP_OK {
		return fmt.Errorf("Error number : %d", retval)
	}
	return nil
}

func (camera *Camera) GetWidgetTree() error {
	var rootWidget *C.CameraWidget

	if retval := C.gp_camera_get_config(camera.gpCamera, (**C.CameraWidget)(unsafe.Pointer(&rootWidget)), camera.gpContext); retval != GP_OK {
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

func (camera *Camera) PrintWidgetTree(file io.Writer) {

	widget := &camera.CameraSettings
	fmt.Printf("Widget Info[%v], Label[%v] , Name[%v], type [%v]\n", widget.Info, widget.Label, widget.Name, widget.Type)
	for _, child := range widget.Children {
		fmt.Printf("    Widget Info[%v], Label[%v] , Name[%v], type [%v]\n", child.Info, child.Label, child.Name, child.Type)
	}

}
