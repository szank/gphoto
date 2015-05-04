package gphoto

// #cgo LDFLAGS: -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include <stdlib.h>
import "C"
import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
)

type cameraFilePathInternal struct {
	Name   [128]uint8
	Folder [1024]uint8
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

	err := getFileBytes(gpFile, buffer)
	if err != nil && leaveOnCamera == false {
		C.gp_camera_file_delete(file.camera.gpCamera, fileDir, fileName, file.camera.gpContext)
	}
	return err
}

//CaptureImage captures image with current setings into camera's internal storage
func (camera *Camera) CaptureImage() (*CameraFilePath, error) {
	photoPath := cameraFilePathInternal{}

	if retval := C.gp_camera_capture(camera.gpCamera, 0, (*C.CameraFilePath)(unsafe.Pointer(&photoPath)), camera.gpContext); retval != gpOk {
		return nil, fmt.Errorf("Cannot capture photo, error code :%v", retval)
	}
	return newCameraFilePathFromInternalImpl(&photoPath, camera), nil
}

func (camera *Camera) CapturePreview(buffer io.Writer) error {
	var gpFile *C.CameraFile
	C.gp_file_new((**C.CameraFile)(unsafe.Pointer(&gpFile)))

	if gpFile == nil {
		return fmt.Errorf("Cannot initialize camera file")
	}
	defer C.gp_file_free(gpFile)

	if retval := C.gp_camera_capture_preview(camera.gpCamera, gpFile, camera.gpContext); retval != gpOk {
		return fmt.Errorf("Cannot capture preview, error code : %d", retval)
	}

	return getFileBytes(gpFile, buffer)
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
func (camera *Camera) ListFiles() ([]CameraStorageInfo, error) {
	var gpCameraStorageInformation *C.CameraStorageInformation
	var storageCount C.int
	storageCount = 0
	returnedStorageInfo := []CameraStorageInfo{}

	retval := C.gp_camera_get_storageinfo(camera.gpCamera, (**C.CameraStorageInformation)(unsafe.Pointer(&gpCameraStorageInformation)), &storageCount, camera.gpContext)
	if retval != gpOk {
		return nil, fmt.Errorf("Cannot get camera storage info, error code %d", retval)
	}
	defer C.free(unsafe.Pointer(gpCameraStorageInformation))

	hdr := reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(gpCameraStorageInformation)),
		Len:  int(storageCount),
		Cap:  int(storageCount),
	}
	nativeCameraFileSystemInfo := *(*[]C.CameraStorageInformation)(unsafe.Pointer(&hdr))
	for i := 0; i < int(storageCount); i++ {
		cameraStorage := CameraStorageInfo{
			Description: C.GoString((*C.char)(&nativeCameraFileSystemInfo[i].description[0])),
			Capacity:    uint64(nativeCameraFileSystemInfo[i].capacitykbytes),
			Free:        uint64(nativeCameraFileSystemInfo[i].freekbytes),
			FreeImages:  uint64(nativeCameraFileSystemInfo[i].freeimages),
			Children:    []CameraFilePath{},

			basedir: C.GoString((*C.char)(&nativeCameraFileSystemInfo[i].basedir[0])),
			camera:  camera,
		}

		if err := camera.recursiveListAllFiles(&cameraStorage.basedir, &cameraStorage.Children); err != nil {
			return nil, err
		}
		returnedStorageInfo = append(returnedStorageInfo, cameraStorage)
	}
	return returnedStorageInfo, nil
}

func (camera *Camera) recursiveListAllFiles(basedir *string, children *[]CameraFilePath) error {
	items, err := camera.findAllChildDirectories(basedir)
	if err != nil {
		return err
	}
	for _, dirName := range items {
		dirItem := CameraFilePath{
			Name:     dirName,
			Folder:   *basedir,
			Isdir:    true,
			Children: []CameraFilePath{},
			camera:   camera,
		}
		childPath := *basedir + "/" + dirName
		if err := camera.recursiveListAllFiles(&childPath, &dirItem.Children); err != nil {
			return err
		}
		*children = append(*children, dirItem)

	}
	items, err = camera.findAllFilesInDir(basedir)
	if err != nil {
		return err
	}
	for _, fileName := range items {
		fileItem := CameraFilePath{
			Name:     fileName,
			Folder:   *basedir,
			Isdir:    false,
			Children: nil,
			camera:   camera,
		}
		*children = append(*children, fileItem)
	}
	fmt.Println("Dir ", *basedir, " has ", len(*children), " children")
	return nil
}

//Hmm, this could be reduced to one func, and a lambda passed as an arg
func (camera *Camera) findAllChildDirectories(basedirPath *string) ([]string, error) {
	var gpFileList *C.CameraList
	var err error
	returnedSlice := []string{}

	gpDirPath := C.CString(*basedirPath)
	defer C.free(unsafe.Pointer(gpDirPath))

	if gpFileList, err = newGphotoList(); err != nil {
		return nil, err
	}
	defer C.gp_list_free(gpFileList)

	if retval := C.gp_camera_folder_list_folders(camera.gpCamera, gpDirPath, gpFileList, camera.gpContext); retval != gpOk {
		return nil, fmt.Errorf("Cannot get folder list from dir %s, error code %v", *basedirPath, retval)
	}

	listSize := int(C.gp_list_count(gpFileList))
	for i := 0; i < listSize; i++ {
		var gpListElementName *C.char
		C.gp_list_get_name(gpFileList, (C.int)(i), (**C.char)(&gpListElementName))
		returnedSlice = append(returnedSlice, C.GoString(gpListElementName))
	}
	return returnedSlice, nil
}

func (camera *Camera) findAllFilesInDir(basedirPath *string) ([]string, error) {
	var gpFileList *C.CameraList
	var err error
	returnedSlice := []string{}

	gpDirPath := C.CString(*basedirPath)
	defer C.free(unsafe.Pointer(gpDirPath))

	if gpFileList, err = newGphotoList(); err != nil {
		return nil, err
	}
	defer C.gp_list_free(gpFileList)

	if retval := C.gp_camera_folder_list_files(camera.gpCamera, gpDirPath, gpFileList, camera.gpContext); retval != gpOk {
		return nil, fmt.Errorf("Cannot get file list from dir %s, error code %v", *basedirPath, retval)
	}

	listSize := int(C.gp_list_count(gpFileList))
	for i := 0; i < listSize; i++ {
		var gpListElementName *C.char
		C.gp_list_get_name(gpFileList, (C.int)(i), (**C.char)(&gpListElementName))
		returnedSlice = append(returnedSlice, C.GoString(gpListElementName))
	}
	return returnedSlice, nil
}
