package gowv

/*
#cgo CFLAGS: -I${SRCDIR}/libs/webview/include
#cgo CXXFLAGS: -I${SRCDIR}/libs/webview/include -DWEBVIEW_STATIC

#cgo linux openbsd freebsd netbsd CXXFLAGS: -DWEBVIEW_GTK -std=c++11
#cgo linux openbsd freebsd netbsd LDFLAGS: -ldl
#cgo linux openbsd freebsd netbsd pkg-config: gtk+-3.0 webkit2gtk-4.1

#cgo darwin CXXFLAGS: -DWEBVIEW_COCOA -std=c++11
#cgo darwin LDFLAGS: -framework WebKit -ldl

#cgo windows CXXFLAGS: -DWEBVIEW_EDGE -std=c++14 -I${SRCDIR}/libs/mswebview2/include
#cgo windows LDFLAGS: -static -ladvapi32 -lole32 -lshell32 -lshlwapi -luser32 -lversion

#include "webview.h"

#include <stdlib.h>
#include <stdbool.h>
#include <stdint.h>

webview_error_t CgoWebViewDispatch(webview_t w, uintptr_t arg);
webview_error_t CgoWebViewBind(webview_t w, const char *name, uintptr_t index);
webview_error_t CgoWebViewUnbind(webview_t w, const char *name);

// NOTE: gowv webview extension c functions.
void CgoNativeWindowSetIcon(void* window_handle, const char* filepath);
void CgoNativeWindowHide(void* window_handle);
void CgoNativeWindowShow(void* window_handle);
// void CgoNativeWindowDecorated(void* window_handle, bool decorated);
void CgoNativeWindowSetFullscreen(void* window_handle, bool fullscreen);
void CgoNativeWindowSetMaximized(void* window_handle);
void CgoNativeWindowSetMinimized(void* window_handle);
*/
import "C"

import (
	"encoding/json"
	"errors"
	"reflect"
	"runtime"
	"sync"
	"unsafe"
)

// Holds the elements of a MAJOR.MINOR.PATCH version number.
type Version struct {
	// Major version.
	Major uint32
	// Minor version.
	Minor uint32
	// Patch version.
	Patch uint32
}

// Holds the library's version information.
type VersionInfo struct {
	// The elements of the version number.
	Version Version
	// SemVer 2.0.0 version number in MAJOR.MINOR.PATCH format.
	VersionNumber string
	// SemVer 2.0.0 pre-release labels prefixed with "-" if specified, otherwise
	// an empty string.
	PreRelease string
	// SemVer 2.0.0 build metadata prefixed with "+", otherwise an empty string.
	BuildMetadata string
}

// Returns the version info from webview used.
func CurrentVersion() VersionInfo {
	version := Version{
		Major: C.WEBVIEW_VERSION_MAJOR,
		Minor: C.WEBVIEW_VERSION_MINOR,
		Patch: C.WEBVIEW_VERSION_PATCH,
	}

	return VersionInfo{
		Version:       version,
		VersionNumber: C.WEBVIEW_VERSION_NUMBER,
		PreRelease:    C.WEBVIEW_VERSION_PRE_RELEASE,
		BuildMetadata: C.WEBVIEW_VERSION_BUILD_METADATA,
	}
}

// Native handle kind. The actual type depends on the backend.
type NativeHandleKind int

const (
	// Top-level window.
	WEBVIEW_NATIVE_HANDLE_KIND_UI_WINDOW NativeHandleKind = iota
	// Browser widget.
	WEBVIEW_NATIVE_HANDLE_KIND_UI_WIDGET
	// Browser controller.
	WEBVIEW_NATIVE_HANDLE_KIND_BROWSER_CONTROLLER
)

// Window size hints.
type Hint int

const (
	// Width and height are default size.
	WEBVIEW_HINT_NONE Hint = iota
	// Width and height are minimum bounds.
	WEBVIEW_HINT_MIN
	// Width and height are maximum bounds.
	WEBVIEW_HINT_MAX
	// Window size can not be changed by a user.
	WEBVIEW_HINT_FIXED
)

// Error codes returned to callers of the API.
type Error int

const (
	// Missing dependency.
	WEBVIEW_ERROR_MISSING_DEPENDENCY Error = -5
	// Operation canceled.
	WEBVIEW_ERROR_CANCELED = -4
	// Invalid state detected
	WEBVIEW_ERROR_INVALID_STATE = -3
	// One or more invalid arguments have been specified e.g. in a function call.
	WEBVIEW_ERROR_INVALID_ARGUMENT = -2
	// An unspecified error occurred. A more specific error code may be needed.
	WEBVIEW_ERROR_UNSPECIFIED = -1
	// OK/Success. Functions that return error codes will typically return this
	// to signify successful operations.
	WEBVIEW_ERROR_OK = 0
	// Signifies that something already exists.
	WEBVIEW_ERROR_DUPLICATE = 1
	// Signifies that something does not exist
	WEBVIEW_ERROR_NOT_FOUND = 2
)

func HadError(err Error) bool {
	return err != WEBVIEW_ERROR_OK
}

func PanicOnError(err Error) {
	if HadError(err) {
		panic(err)
	}
}

// Webview instance wrapper.
type Instance struct {
	// Pointer to a webver instance.
	W C.webview_t
}

// Global vars.
var (
	mu       sync.Mutex
	index    uintptr
	dispatch = map[uintptr]func(){}
	bindings = map[uintptr]func(id string, req string) (any, error){}
)

func init() {
	// native GUI toolkits require main OS thread
	runtime.LockOSThread()
}

//export _webviewDispatchGoCallback
func _webviewDispatchGoCallback(index unsafe.Pointer) {
	mu.Lock()
	f := dispatch[uintptr(index)]
	delete(dispatch, uintptr(index))
	mu.Unlock()
	f()
}

//export _webviewBindingGoCallback
func _webviewBindingGoCallback(w C.webview_t, id *C.char, req *C.char, index uintptr) {
	mu.Lock()
	f := bindings[uintptr(index)]
	mu.Unlock()
	jsString := func(v interface{}) string { b, _ := json.Marshal(v); return string(b) }
	status, result := 0, ""
	if res, err := f(C.GoString(id), C.GoString(req)); err != nil {
		status = -1
		result = jsString(err.Error())
	} else if b, err := json.Marshal(res); err != nil {
		status = -1
		result = jsString(err.Error())
	} else {
		status = 0
		result = string(b)
	}
	s := C.CString(result)
	defer C.free(unsafe.Pointer(s))

	PanicOnError(Error(C.webview_return(w, id, C.int(status), s)))
}

// Creates a new webview instance.
func (h *Instance) Create(debug bool, window unsafe.Pointer) {
	var d int
	if debug {
		d = 1
	} else {
		d = 0
	}

	h.W = C.webview_create(C.int(d), window)
}

// Destroys a webview instance and closes the native window.
func (h *Instance) Destroy() Error {
	return Error(C.webview_destroy(h.W))
}

// Runs the main loop until it's terminated.
func (h *Instance) Run() Error {
	return Error(C.webview_run(h.W))
}

// Stops the main loop. It is safe to call this function from another
// background thread.
func (h *Instance) Terminate() Error {
	mu.Lock()
	err := C.webview_terminate(h.W)
	mu.Unlock()

	return Error(err)
}

func (h *Instance) GetWindow() unsafe.Pointer {
	return C.webview_get_window(h.W)
}

func (h *Instance) GetNativeHandle(kind NativeHandleKind) unsafe.Pointer {
	return C.webview_get_native_handle(h.W, C.webview_native_handle_kind_t(kind))
}

// Schedules a function to be invoked on the thread with the run/event loop.
//
// Since library functions generally do not have thread safety guarantees,
// this function can be used to schedule code to execute on the main/GUI
// thread and thereby make that execution safe in multi-threaded applications.
// TODO
func (h *Instance) Dispatch(fn func()) Error {
	mu.Lock()
	for ; dispatch[index] != nil; index++ {
	}
	dispatch[index] = fn
	mu.Unlock()

	return Error(C.CgoWebViewDispatch(h.W, C.uintptr_t(index)))
}

// Updates the title of the native window.
func (h *Instance) SetTitle(title string) Error {
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))

	return Error(C.webview_set_title(h.W, s))
}

// Updates the size of the native window
func (h *Instance) SetSize(width int, height int, hints Hint) Error {
	return Error(C.webview_set_size(h.W, C.int(width), C.int(height), C.webview_hint_t(hints)))
}

// Set the runtime icon of the native window.
func (h *Instance) SetIcon(icon string) Error {
	s := C.CString(icon)
	defer C.free(unsafe.Pointer(s))

	C.CgoNativeWindowSetIcon(h.GetWindow(), s)

	return WEBVIEW_ERROR_OK
}

// Hides the current window.
func (h *Instance) Hide() Error {
	C.CgoNativeWindowHide(h.GetWindow())

	return WEBVIEW_ERROR_OK
}

// Shows the current window.
func (h *Instance) Show() Error {
	C.CgoNativeWindowShow(h.GetWindow())

	return WEBVIEW_ERROR_OK
}

// Makes window fullscreen.
func (h *Instance) SetFullscreen(fullscreen bool) Error {
	C.CgoNativeWindowSetFullscreen(h.GetWindow(), C.bool(fullscreen))

	return WEBVIEW_ERROR_OK
}

// Sets window to maximum bounds.
func (h *Instance) SetMaximized() Error {
	C.CgoNativeWindowSetMaximized(h.GetWindow())

	return WEBVIEW_ERROR_OK
}

// Iconifies the window.
func (h *Instance) SetMinimized() Error {
	C.CgoNativeWindowSetMinimized(h.GetWindow())

	return WEBVIEW_ERROR_OK
}

// Navigates webview to the given URL. URL may be a properly encoded data URI.
func (h *Instance) Navigate(url string) Error {
	s := C.CString(url)
	defer C.free(unsafe.Pointer(s))

	return Error(C.webview_navigate(h.W, s))
}

// Load HTML content into the webview.
func (h *Instance) SetHTML(html string) Error {
	s := C.CString(html)
	defer C.free(unsafe.Pointer(s))

	return Error(C.webview_set_html(h.W, s))
}

// Injects JavaScript code to be executed immediately upon loading a page.
// The code will be executed before window.onload.
func (h *Instance) Init(js string) Error {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))

	return Error(C.webview_init(h.W, s))
}

// Evaluates arbitrary JavaScript code.
func (h *Instance) Eval(js string) Error {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))

	return Error(C.webview_eval(h.W, s))
}

// Binds a function pointer to a new global JavaScript function.
// TODO
func (h *Instance) Bind(name string, fn any) Error {
	v := reflect.ValueOf(fn)

	if v.Kind() != reflect.Func {
		return WEBVIEW_ERROR_UNSPECIFIED
	}

	if n := v.Type().NumOut(); n > 2 {
		return WEBVIEW_ERROR_INVALID_ARGUMENT
	}

	binding := func(id, req string) (interface{}, error) {
		raw := []json.RawMessage{}
		if err := json.Unmarshal([]byte(req), &raw); err != nil {
			return nil, err
		}

		isVariadic := v.Type().IsVariadic()
		numIn := v.Type().NumIn()
		if (isVariadic && len(raw) < numIn-1) || (!isVariadic && len(raw) != numIn) {
			return nil, errors.New("function arguments mismatch")
		}
		args := []reflect.Value{}
		for i := range raw {
			var arg reflect.Value
			if isVariadic && i >= numIn-1 {
				arg = reflect.New(v.Type().In(numIn - 1).Elem())
			} else {
				arg = reflect.New(v.Type().In(i))
			}
			if err := json.Unmarshal(raw[i], arg.Interface()); err != nil {
				return nil, err
			}
			args = append(args, arg.Elem())
		}
		errorType := reflect.TypeOf((*error)(nil)).Elem()
		res := v.Call(args)
		switch len(res) {
		case 0:
			return nil, nil
		case 1:
			if res[0].Type().Implements(errorType) {
				if res[0].Interface() != nil {
					return nil, res[0].Interface().(error)
				}
				return nil, nil
			}
			return res[0].Interface(), nil
		case 2:
			if !res[1].Type().Implements(errorType) {
				return nil, errors.New("second return value must be an error")
			}
			if res[1].Interface() == nil {
				return res[0].Interface(), nil
			}
			return res[0].Interface(), res[1].Interface().(error)
		default:
			return nil, errors.New("unexpected number of return values")
		}
	}

	mu.Lock()
	for ; bindings[index] != nil; index++ {
	}
	bindings[index] = binding
	mu.Unlock()
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return Error(C.CgoWebViewBind(h.W, cname, C.uintptr_t(index)))
}

// Removes a binding created with [Instance].Bind.
// TODO
func (h *Instance) Unbind(name string) Error {
	cname := C.CString(name)
	defer C.free(unsafe.Pointer(cname))

	return Error(C.CgoWebViewUnbind(h.W, cname))
}

// Responds to a binding call from the JS side.
func (h *Instance) Return(id string, status int, result string) Error {
	i := C.CString(id)
	defer C.free(unsafe.Pointer(i))

	r := C.CString(id)
	defer C.free(unsafe.Pointer(r))

	return Error(C.webview_return(h.W, i, C.int(status), r))
}
