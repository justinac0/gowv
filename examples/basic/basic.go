package main

import (
	"gowv"
)

func main() {
	w := gowv.Instance{}

	w.Create(true, nil)

	defer func() {
		gowv.PanicOnError(w.Destroy())
	}()

	gowv.PanicOnError(w.SetTitle("Basic Example"))
	gowv.PanicOnError(w.SetSize(480, 320, gowv.WEBVIEW_HINT_NONE))
	gowv.PanicOnError(w.SetHTML("Thanks for using webview"))
	gowv.PanicOnError(w.Run())
}
