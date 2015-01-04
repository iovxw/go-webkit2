package webkit2

// #include <gio/gio.h>
// #include "gasyncreadycallback.go.h"
import "C"
import (
	"errors"
	"fmt"
	"reflect"
	"unsafe"
)

type garCallback struct {
	f reflect.Value
}

//export _go_gasyncreadycallback_call
func _go_gasyncreadycallback_call(cbinfoRaw C.gpointer, cresult unsafe.Pointer) {
	result := (*C.GAsyncResult)(cresult)
	fmt.Println("cresult:%v", cresult)
	fmt.Println("result:%v", result)
	cbinfo := (*garCallback)(unsafe.Pointer(cbinfoRaw))
	fmt.Println("cbinfo:%v", cbinfo)

	if result == nil {
		cbinfo.f.Call(nil)
		return
	}
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in f", r)
			//cbinfo.f.Call(nil)
		}
	}()
	cbinfo.f.Call([]reflect.Value{reflect.ValueOf(result)})
}

func newGAsyncReadyCallback(f interface{}) (cCallback C.GAsyncReadyCallback, userData C.gpointer, err error) {
	rf := reflect.ValueOf(f)
	if rf.Kind() != reflect.Func {
		return nil, nil, errors.New("f is not a function")
	}
	cbinfo := &garCallback{rf}
	return C.GAsyncReadyCallback(C._gasyncreadycallback_call), C.gpointer(unsafe.Pointer(cbinfo)), nil
}
