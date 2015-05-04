package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include <stdlib.h>
import "C"
import (
	"bytes"
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

func newCameraFilePathFromInternalImpl(input *cameraFilePathInternal, camera *Camera) *CameraFilePath {
	return &CameraFilePath{
		Name:     string(input.Name[:bytes.IndexByte(input.Name[:], 0)]),
		Folder:   string(input.Folder[:bytes.IndexByte(input.Folder[:], 0)]),
		Isdir:    false,
		Children: nil,
		camera:   camera,
	}
}

func newGphotoList() (*C.CameraList, error) {
	var gpFileList *C.CameraList
	if retval := C.gp_list_new((**C.CameraList)(unsafe.Pointer(&gpFileList))); retval != gpOk {
		return nil, fmt.Errorf("Could not create a list, eror code %v", retval)
	}
	return gpFileList, nil
}

func getFileBytes(gpFileIn *C.CameraFile, bufferOut io.Writer) error {
	var fileData *C.char
	var fileLen C.ulong
	C.gp_file_get_data_and_size(gpFileIn, (**C.char)(unsafe.Pointer(&fileData)), &fileLen)

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(fileData)),
		Len:  int(fileLen),
		Cap:  int(fileLen),
	}
	goSlice := *(*[]byte)(unsafe.Pointer(&hdr))

	_, err := bufferOut.Write(goSlice)
	return err
}
