package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
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

type cameraFilePathInternal struct {
	Name   [128]uint8
	Folder [1024]uint8
}

//CameraFilePath is a path to a file or dir on the camera file system
type CameraFilePath struct {
	Name     string
	Folder   string
	Isdir    bool
	Children []CameraFilePath
	camera   *Camera
}

//DownloadImage saves image pointed by path to the provided buffer. If leave on camera is set to false,the file will be deleted from the camera internal storage
func (file *CameraFilePath) DownloadImage(buffer io.Writer, leaveOnCamera bool) error {
	var gpFile *C.CameraFile
	C.gp_file_new((**C.CameraFile)(unsafe.Pointer(&gpFile)))

	if gpFile == nil {
		return fmt.Errorf("Cannot initialize camera file")
	}
	defer C.gp_file_free(gpFile)

	fileDir := C.CString(file.Folder)
	defer C.free(unsafe.Pointer(fileDir))

	fileName := C.CString(file.Name)
	defer C.free(unsafe.Pointer(fileName))

	if retval := C.gp_camera_file_get(file.camera.gpCamera, fileDir, fileName, FileTypeNormal, gpFile, file.camera.gpContext); retval != gpOk {
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
		C.gp_camera_file_delete(file.camera.gpCamera, fileDir, fileName, file.camera.gpContext)
	}
	return err
}

func newCameraFilePathFromInternalImpl(input *cameraFilePathInternal, camera *Camera) *CameraFilePath {
	return &CameraFilePath{
		Name:     string(input.Name[:bytes.IndexByte(input.Name[:], 0)]),
		Folder:   string(input.Folder[:bytes.IndexByte(input.Folder[:], 0)]),
		Isdir:    false,
		Children: nil,
		camera:   camera,
	}
}

//CaptureImage captures image with current setings into camera's internal storage
func (camera *Camera) CaptureImage() (*CameraFilePath, error) {
	photoPath := cameraFilePathInternal{}

	if retval := C.gp_camera_capture(camera.gpCamera, 0, (*C.CameraFilePath)(unsafe.Pointer(&photoPath)), camera.gpContext); retval != gpOk {
		return nil, fmt.Errorf("Cannot capture photo, error code :%v", retval)
	}
	return newCameraFilePathFromInternalImpl(&photoPath, camera), nil
}

//DeleteFile tries to delete file from the camera, and returns error if it fails
func (camera *Camera) DeleteFile(path *CameraFilePath) error {
	fileDir := C.CString(path.Folder)
	defer C.free(unsafe.Pointer(fileDir))

	fileName := C.CString(path.Name)
	defer C.free(unsafe.Pointer(fileName))
	retval := C.gp_camera_file_delete(camera.gpCamera, fileDir, fileName, camera.gpContext)
	if retval < 0 {
		return fmt.Errorf("Cannot delete fine on camera, error code :%v", retval)
	}
	return nil
}

//ListFiles returns a lits of files and folders on the camera
//Currently work in progress
func (camera *Camera) ListFiles() ([]CameraFilePath, error) {
	return nil, nil
}
