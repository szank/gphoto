package gphoto

// #cgo LDFLAGS: -L. -lgphotocallbacks -L/usr/lib/x86_64-linux-gnu -lgphoto2 -lgphoto2_port
// #cgo CFLAGS: -I/usr/include
// #include <gphoto2/gphoto2.h>
// #include "callbacks.h"
import "C"
import "fmt"

//Context represents a  context in which all other calls are executed
type Context struct {
	gpContext *C.GPContext
}

//Free should be called afer you don't need the context anymore
func (c Context) Free() {
	C.gp_context_unref(c.gpContext)
}

//GetNewGPhotoContext returns a new gphoto context
func GetNewGPhotoContext() (*Context, error) {
	var gpContext *C.GPContext
	fmt.Printf("Gpcontext before call %#v\n", gpContext)
	gpContext = C.gp_context_new()
	fmt.Printf("Gpcontext after call %#v\n", gpContext)

	if gpContext == nil {
		return nil, fmt.Errorf("Could not initialize libgphoto2 context")
	}

	C.gp_context_set_error_func(gpContext, (*[0]byte)(C.ctx_error_func), nil)
	C.gp_context_set_status_func(gpContext, (*[0]byte)(C.ctx_status_func), nil)
	return &Context{
		gpContext: gpContext,
	}, nil
}
