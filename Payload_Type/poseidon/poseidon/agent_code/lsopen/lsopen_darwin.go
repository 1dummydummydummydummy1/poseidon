// +build darwin

package lsopen

/*
#cgo CFLAGS: -x objective-c -fmacro-backtrace-limit=0 -std=gnu11 -Wobjc-property-no-attribute -Wunguarded-availability-new
#cgo LDFLAGS: -framework Foundation -framework CoreServices
#include "lsopen_darwin.h"
#include <stdlib.h>
*/
import "C"
import "unsafe"
import "os"

type DyldInjectDarwin struct {
	Successful bool
}

func (j *DyldInjectDarwin) Success() bool {
	return j.Successful
}

func runCommand(app string, hide bool, args []string) (DyldInjectDarwin, error) {
	capp := C.CString(app)

	var c_argc C.int = 0
	var c_argv **C.char = nil
	var chide int

	if hide {
		chide = 1
	} else {
		chide = 0
	}

	ihide := C.int(chide)

	//prepping args to pass to function
	c_argc = C.int(len(args) + 1)
	cArgs := make([](*C.char), len(args)+2)
	for i := range cArgs {
		cArgs[i] = nil
	}
	cArgs[0] = C.CString(os.Args[0])
	for i, arg := range args {
		cArgs[i+1] = C.CString(arg)
	}
	c_argv = (**C.char)(unsafe.Pointer(&cArgs[0]))


	res := C.dyldd_inject(capp, ihide, c_argv, c_argc)

	//free
	for i := range cArgs {
		if cArgs[i] != nil {
			defer C.free(unsafe.Pointer(cArgs[i]))
		}
	}

	r := DyldInjectDarwin{}
	if res == 0 {
		r.Successful = true
	} else {
		r.Successful = false
	}

	return r, nil
} //hiiii
