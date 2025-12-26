package gowv_test

import (
	"testing"

	"github.com/justinac0/gowv"
)

func TestWebView(t *testing.T) {
	w := gowv.Instance{}

	w.Create(true, nil)

	defer func() {
		gowv.PanicOnError(w.Destroy())
	}()

	gowv.PanicOnError(w.SetTitle("Testing Webview"))
	gowv.PanicOnError(w.SetSize(480, 320, gowv.WEBVIEW_HINT_NONE))
	gowv.PanicOnError(w.SetHTML("Close window to end test"))
	gowv.PanicOnError(w.Run())
}
