package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"unsafe"
)

//Camera struct represents a camera connected to the computer
type Camera struct {
	gpCamera       *C.Camera
	gpContext      *C.GPContext
	CameraSettings CameraWidget
}

//CameraWidget is a representation of one of the cameras's setting or control knob
type CameraWidget interface {
	Label() string
	Name() string
	Info() string
	Type() WidgetType
	Children() []CameraWidget
	ReadOnly() bool
}

type cameraWidgetImpl struct {
	label      string
	name       string
	info       string
	widgetType WidgetType
	children   []CameraWidget
	readonly   bool
	camera     *Camera
}

func (w cameraWidgetImpl) Label() string {
	return w.label
}

func (w cameraWidgetImpl) Name() string {
	return w.name
}

func (w cameraWidgetImpl) Info() string {
	return w.info
}

func (w cameraWidgetImpl) Type() WidgetType {
	return w.widgetType
}

func (w cameraWidgetImpl) Children() []CameraWidget {
	return w.children
}

func (w cameraWidgetImpl) ReadOnly() bool {
	return w.readonly
}

type CameraWidgetText interface {
	Get() (*string, error)
	Set(*string) error
}

type cameraWidgetTextImpl struct {
	cameraWidgetImpl
}

func (w cameraWidgetTextImpl) Get() (*string, error) {
	var gpText *C.char
	var gpWidget *C.CameraWidget
	var err error
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return nil, err
	}
	defer w.camera.freeChildWidget(gpWidget)

	if retval := C.gp_widget_get_value(gpWidget, (unsafe.Pointer(&gpText))); retval != gpOk {
		return nil, fmt.Errorf("Cannot read widget property value, error code :%d", retval)
	}
	widgetValue := C.GoString(gpText)
	return &widgetValue, nil
}

func (w cameraWidgetTextImpl) Set(input *string) error {
	var err error
	gpText := C.CString(*input)
	defer C.free(unsafe.Pointer(gpText))

	var gpWidget *C.CameraWidget
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return err
	}
	defer w.camera.freeChildWidget(gpWidget)

	if retval := C.gp_widget_set_value(gpWidget, unsafe.Pointer(gpText)); retval != gpOk {
		return fmt.Errorf("Could not set widget value , error code :%d", retval)
	}
	return nil
}

type CameraWidgetMenu interface {
	CameraWidgetText
	GetChoices() ([]string, error)
}

type cameraWidgetMenuImpl struct {
	cameraWidgetTextImpl
}

func (w cameraWidgetMenuImpl) Set(input *string) error {
	if choices, err := w.GetChoices(); err != nil {
		for _, item := range choices {
			if item == *input {
				return w.cameraWidgetTextImpl.Set(input)
			}
		}
		return fmt.Errorf("Could not find provided value in alloved values list")
	} else {
		return err
	}
}

func (w cameraWidgetMenuImpl) GetChoices() ([]string, error) {
	var gpWidget *C.CameraWidget
	var err error
	if gpWidget, err = w.camera.getChildWidget(&w.name); err != nil {
		return nil, err
	}
	defer w.camera.freeChildWidget(gpWidget)

	choicesList := []string{}
	numChoices := C.gp_widget_count_choices(gpWidget)
	for i := 0; i < int(numChoices); i++ {
		var gpChoice *C.char
		C.gp_widget_get_choice(gpWidget, C.int(i), (**C.char)(unsafe.Pointer(&gpChoice)))
		choicesList = append(choicesList, C.GoString(gpChoice))
	}
	return choicesList, nil
}

type CameraWidgetRadio CameraWidgetMenu

type WidgetType string

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

//CamersStorageInfo is a struct describing one of the camera's storage spaces (SD or CF cards for example)
//Children is a directory tree present on the storage space
type CameraStorageInfo struct {
	Description string
	Capacity    uint64
	Free        uint64
	FreeImages  uint64
	Children    []CameraFilePath

	basedir string
	camera  *Camera
}
