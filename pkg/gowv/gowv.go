package gowv

/*
#cgo CFLAGS: -I${SRCDIR}/libs/webview/include
#cgo CXXFLAGS: -I${SRCDIR}/libs/webview/include -DWEBVIEW_STATIC -DWEBVIEW_IMPLEMENTATION

#cgo linux CXXFLAGS: -DWEBVIEW_GTK -std=c++11
#cgo linux LDFLAGS: -ldl
#cgo linux pkg-config: gtk+-3.0 webkit2gtk-4.1

#include "webview.h"
#include <stdlib.h>
#include <stdint.h>
*/
import "C"

import (
	"fmt"
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
		VersionNumber: fmt.Sprintf("%v.%v.%v", version.Major, version.Minor, version.Patch),
		PreRelease:    fmt.Sprintf("-%v.%v.%v", version.Major, version.Minor, version.Patch),
		BuildMetadata: fmt.Sprintf("+%v.%v.%v", version.Major, version.Minor, version.Patch),
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
	UWEBVIEW_ERROR_NSPECIFIED = -1
	// OK/Success. Functions that return error codes will typically return this
	// to signify successful operations.
	WEBVIEW_ERROR_OK = 0
	// Signifies that something already exists.
	WEBVIEW_ERROR_DUPLICATE = 1
	// Signifies that something does not exist
	WEBVIEW_ERROR_NOT_FOUND = 2
)

func HadError(err Error) bool {
	return err == WEBVIEW_ERROR_OK
}

// Webview instance wrapper.
type Instance struct {
	// Pointer to a webver instance.
	W C.webview_t
}

// Global vars.
var (
	mu         sync.Mutex
	dispatches map[string]func()
	functions  map[string]func()
)

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

// Schedules a function to be invoked on the thread with the run/event loop.
//
// Since library functions generally do not have thread safety guarantees,
// this function can be used to schedule code to execute on the main/GUI
// thread and thereby make that execution safe in multi-threaded applications.
func (h *Instance) Dispatch() Error {

	return WEBVIEW_ERROR_OK
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

// Set the icon of the native window
func (h *Instance) SetIcon(icon string) Error {
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

type BindFun func(id string, req string, arg unsafe.Pointer) 

// Binds a function pointer to a new global JavaScript function.
func (h *Instance) Bind(name string, fn BindFun, arg unsafe.Pointer) Error {
	fmt.Println("binding: ", name)

	// check if function exists in map
	// map function

	return WEBVIEW_ERROR_OK
}

// Removes a binding created with [Instance].Bind.
func (h *Instance) Unbind(name string) Error {
	fmt.Println("unbind: ", name)

	s := C.CString(name)
	defer C.free(unsafe.Pointer(s))

	// check if function exists in map
	// remove bound function from map
	err := C.webview_unbind(h.W, s)

	return Error(err)
}

// Responds to a binding call from the JS side.
func (h *Instance) Return(id string, status int, result string) Error {
	// This function is safe to call from another thread.
	mu.Lock()
	defer mu.Unlock()

	i := C.CString(id)
	defer C.free(unsafe.Pointer(i))

	r := C.CString(id)
	defer C.free(unsafe.Pointer(r))

	return Error(C.webview_return(h.W, i, C.int(status), r))
}
