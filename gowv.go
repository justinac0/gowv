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

type handle struct {
	W C.webview_t
}

type WebView interface {
	Create(debug int, window unsafe.Pointer)
	Destroy() Error
	SetTitle(title string) Error
	SetSize(width int, height int, hints Hint) Error
	SetHTML(html string) Error
	Run() Error
}

func (h *handle) Create(debug bool, window unsafe.Pointer) {
	var d int
	if debug {
		d = 1
	} else {
		d = 0
	}

	h.W = C.webview_create(C.int(d), window)
}

func (h *handle) Destroy() Error {
	err := C.webview_destroy(h.W)
	fmt.Println("error", err)

	return WEBVIEW_ERROR_OK
}

func (h *handle) SetTitle(title string) Error {
	s := C.CString(title)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_set_title(h.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func (h *handle) SetSize(width int, height int, hints Hint) Error {
	err := C.webview_set_size(h.W, C.int(width), C.int(height), C.webview_hint_t(hints))
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func (h *handle) SetHTML(html string) Error {
	s := C.CString(html)
	defer C.free(unsafe.Pointer(s))

	err := C.webview_set_html(h.W, s)
	fmt.Println("error ", err)

	return WEBVIEW_ERROR_OK
}

func (h *handle) Run() Error {
	err := C.webview_run(h.W)
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

	w := handle{}
	w.Create(false, nil)
	defer w.Destroy()

	w.SetTitle("Hello, World")
	w.SetSize(640, 480, WEBVIEW_HINT_NONE)
	w.SetHTML("<p>Hello, World</p>")
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
