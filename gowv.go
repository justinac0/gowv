package main

/*
#cgo CFLAGS: -I./libs/webview/include
#cgo CXXFLAGS: -I./libs/webview/include -DWEBVIEW_STATIC -DWEBVIEW_IMPLEMENTATION

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
	"runtime"
	"unsafe"
)

type NativeHandleKind int

const (
	WEBVIEW_NATIVE_HANDLE_KIND_UI_WINDOW NativeHandleKind = iota
	WEBVIEW_NATIVE_HANDLE_KIND_UI_WIDGET
	WEBVIEW_NATIVE_HANDLE_KIND_BROWSER_CONTROLLER
)

type Hint int

const (
	WEBVIEW_HINT_NONE Hint = iota
	WEBVIEW_HINT_MIN
	WEBVIEW_HINT_MAX
	WEBVIEW_HINT_FIXED
)

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

type Version struct {
	Major uint32
	Minor uint32
	Patch uint32
}

type VersionInfo struct {
	Version       Version
	VersionNumber [32]byte
	PreRelease    [48]byte
	BuildMetadata [48]byte
}

type webview struct {
	W C.webview_t
}

// Creates a new webview instance.
func (w *webview) Create(debug bool, window unsafe.Pointer) {
	var d int
	if debug {
		d = 1
	} else {
		d = 0
	}

	w.W = C.webview_create(C.int(d), window)
}

// Destroys a webview instance and closes the native window.
func (w *webview) Destroy() Error {
	err := C.webview_destroy(w.W)
	fmt.Println("error", err)

	return WEBVIEW_ERROR_OK
}

// Runs the main loop until it's terminated.
func (w *webview) Run() Error {
	err := C.webview_run(w.W)
	fmt.Println("error", err)

	return WEBVIEW_ERROR_OK
}

func (w *webview) SetTitle(title string) Error {
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_set_title(w.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func (w *webview) SetSize(width int, height int, hints Hint) Error {
	err := C.webview_set_size(w.W, C.int(width), C.int(height), C.webview_hint_t(hints))
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

// Navigates webview to the given URL. URL may be a properly encoded data URI.
func (w *webview) Navigate(url string) Error {
	s := C.CString(url)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_navigate(w.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

// Load HTML content into the webview.
func (w *webview) SetHTML(html string) Error {
	s := C.CString(html)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_set_html(w.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func (w *webview) Init(js string) Error {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_init(w.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func (w *webview) Eval(js string) Error {
	s := C.CString(js)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_eval(w.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func CurrentVersion() VersionInfo {
	version := Version{
		Major: C.WEBVIEW_VERSION_MAJOR,
		Minor: C.WEBVIEW_VERSION_MINOR,
		Patch: C.WEBVIEW_VERSION_PATCH,
	}

	return VersionInfo{
		Version: version,
	}
}

func main() {
	runtime.LockOSThread()

	w := webview{}
	w.Create(true, nil)
	defer w.Destroy()

	w.SetTitle("Hello, World")
	w.SetSize(640, 480, WEBVIEW_HINT_NONE)
	w.Navigate("http://192.168.0.236:3000")
	w.Run()

	v := CurrentVersion()
	fmt.Println(v.Version.Major, v.Version.Minor, v.Version.Patch)
}

/*
  webview_t w = webview_create(0, NULL);
  webview_set_title(w, "Basic Example");
  webview_set_size(w, 480, 320, WEBVIEW_HINT_NONE);
  webview_set_html(w, "Thanks for using webview!");
  webview_run(w);
  webview_destroy(w);
*/
