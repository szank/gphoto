package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

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
	var err error

	if rootWidget, err = camera.getRootWidget(); err != nil {
		return err
	}
	defer C.free(unsafe.Pointer(rootWidget))

	camera.CameraSettings = camera.getWidgetInfo(rootWidget)
	return nil
}

func (camera *Camera) getWidgetInfo(input *C.CameraWidget) CameraWidget {
	//	var child *C.CameraWidget
	var gpInfo *C.char
	var gpLabel *C.char
	var gpName *C.char
	var gpWidgetType C.CameraWidgetType
	var child *C.CameraWidget
	var readonly C.int

	//ignore return value (int) for now
	C.gp_widget_get_info(input, (**C.char)(unsafe.Pointer(&gpInfo)))
	C.gp_widget_get_label(input, (**C.char)(unsafe.Pointer(&gpLabel)))
	C.gp_widget_get_name(input, (**C.char)(unsafe.Pointer(&gpName)))
	C.gp_widget_get_type(input, (*C.CameraWidgetType)(unsafe.Pointer(&gpWidgetType)))
	C.gp_widget_get_readonly(input, &readonly)

	baseWidget := cameraWidgetImpl{
		widgetType: widgetType(gpWidgetType),
		label:      C.GoString(gpLabel),
		info:       C.GoString(gpInfo),
		name:       C.GoString(gpName),
		readonly:   (int(readonly) == 1),
		camera:     camera,
	}

	childrenCount := int(C.gp_widget_count_children(input))
	for n := 0; n < childrenCount; n++ {
		C.gp_widget_get_child(input, C.int(n), (**C.CameraWidget)(unsafe.Pointer(&child)))
		baseWidget.children = append(baseWidget.children, camera.getWidgetInfo(child))
	}

	switch baseWidget.widgetType {
	case WidgetText:
		return cameraWidgetTextImpl{
			cameraWidgetImpl: baseWidget,
		}
	case WidgetMenu, WidgetRadio:
		return cameraWidgetMenuImpl{
			cameraWidgetTextImpl: cameraWidgetTextImpl{
				cameraWidgetImpl: baseWidget,
			},
		}
	default:
		return baseWidget
	}
}

//PrintWidgetTree prints widget hierarchy to the output buffer
func (camera *Camera) PrintWidgetTree(file io.Writer) {
	camera.printWidgetTreeRecursive(file, camera.CameraSettings)
}

func (camera *Camera) printWidgetTreeRecursive(file io.Writer, widget CameraWidget) {
	fmt.Printf("Widget Info[%v], Label[%v] , Name[%v], readonly [%t], type [%v]\n", widget.Info(), widget.Label(), widget.Name(), widget.ReadOnly(), widget.Type())
	switch widgetAccessor := widget.(type) {
	case CameraWidgetMenu:
		if choices, err := widgetAccessor.GetChoices(); err != nil {
			fmt.Println("     Could not read value of the widget")
		} else {
			fmt.Printf("      Choices : %+v", choices)
		}
		if value, err := widgetAccessor.Get(); err != nil {
			fmt.Println("     Could not read value of the widget")
		} else {
			fmt.Printf("      Value :%s\n", *value)
		}
	case CameraWidgetRadio:
		if choices, err := widgetAccessor.GetChoices(); err != nil {
			fmt.Println("     Could not read value of the widget")
		} else {
			fmt.Printf("      Choices : %+v", choices)
		}
		if value, err := widgetAccessor.Get(); err != nil {
			fmt.Println("     Could not read value of the widget")
		} else {
			fmt.Printf("      Value :%s\n", *value)
		}
	case CameraWidgetText:
		if value, err := widgetAccessor.Get(); err != nil {
			fmt.Println("     Could not read value of the widget")
		} else {
			fmt.Printf("      Value :%s\n", *value)
		}

	default:
		fmt.Printf("     Cannot get value for the type %v\n", reflect.TypeOf(widgetAccessor))

	}
	for _, child := range widget.Children() {
		camera.printWidgetTreeRecursive(file, child)
	}
}

func (camera *Camera) getRootWidget() (*C.CameraWidget, error) {
	var rootWidget *C.CameraWidget

	if retval := C.gp_camera_get_config(camera.gpCamera, (**C.CameraWidget)(unsafe.Pointer(&rootWidget)), camera.gpContext); retval != gpOk {
		return nil, fmt.Errorf("cannot initialize camera settings tree error code :%v", retval)
	}
	return rootWidget, nil
}

func (camera *Camera) getChildWidget(name *string) (*C.CameraWidget, error) {
	var rootWidget, childWidget *C.CameraWidget
	var err error
	if rootWidget, err = camera.getRootWidget(); err != nil {
		return nil, err
	}

	gpChildWidgetName := C.CString(*name)
	defer C.free(unsafe.Pointer(gpChildWidgetName))

	if retval := C.gp_widget_get_child_by_name(rootWidget, gpChildWidgetName, (**C.CameraWidget)(unsafe.Pointer(&childWidget))); retval != gpOk {
		return nil, fmt.Errorf("Could not retrieve child widget with name %s, error code %d", *name, retval)
	}
	return childWidget, nil
}

func (camera *Camera) freeChildWidget(input *C.CameraWidget) {
	var rootWidget *C.CameraWidget
	C.gp_widget_get_root(input, (**C.CameraWidget)(unsafe.Pointer(&rootWidget)))
	C.free(unsafe.Pointer(rootWidget))
}
