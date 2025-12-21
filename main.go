package main

import (
	"fmt"
	gowv "gowv/pkg/gowv"
	"unsafe"
)

type Context struct {
	count int64 
}

type BindFun func(id string, req string, arg unsafe.Pointer) 
func count(id string, req string, arg unsafe.Pointer) {
	fmt.Println("hello")
}

func main() {
	w := gowv.Instance{}
	w.Create(true, nil)

	defer w.Destroy()

	c := Context{}

	w.SetTitle("Hello, World")
	w.SetSize(1280, 720, gowv.WEBVIEW_HINT_NONE)
	w.Navigate("http://192.168.0.236:3000")
	w.Unbind("count")
	w.Bind("count", count, unsafe.Pointer(&c))
	w.Run()
}
