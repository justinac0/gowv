// package main

// import (
// 	"gowv"
// )

// func main() {
// 	w := gowv.Instance{}

// 	w.Create(true, nil)

// 	defer func() {
// 		gowv.PanicOnError(w.Destroy())
// 	}()

// 	gowv.PanicOnError(w.SetTitle("Basic Example"))
// 	gowv.PanicOnError(w.SetSize(480, 320, gowv.WEBVIEW_HINT_NONE))
// 	gowv.PanicOnError(w.SetHTML("Thanks for using webview"))
// 	w.SetIcon("hello")
// 	w.Show()
// 	w.Hide()
// 	w.SetMaximized()
// 	w.SetMinimized()
// 	gowv.PanicOnError(w.Run())
// }


package main

import webview "gowv"

const html = `
<button id="set_icon">Set Icon</button>
<button id="show">Show</button>
<button id="hide">Hide</button>
<button id="set_maximized">Set Maximized</button>
<button id="set_minimized">Set Minimized</button>
<script>
  const setIcon = document.querySelector("#set_icon")
  const show = document.querySelector("#show")
  const hide = document.querySelector("#hide")
  const set_maximized = document.querySelector("#set_maximized")
  const set_minimized = document.querySelector("#set_minimized")
  document.addEventListener("DOMContentLoaded", () => {
    setIcon.addEventListener("click", () => { window.native_window_set_icon().then(result => {}); });
    show.addEventListener("click", () => { window.native_window_show().then(result => {}); });
    hide.addEventListener("click", () => { window.native_window_hide().then(result => {}); });
    set_maximized.addEventListener("click", () => { window.native_window_set_maximized().then(result => {}); });
    set_minimized.addEventListener("click", () => { window.native_window_set_minimized().then(result => {}); });
  });
</script>`

func main() {
	w := webview.Instance{}
	w.Create(true, nil)
	defer w.Destroy()
	w.SetTitle("Extended Example")
	w.SetSize(480, 320, webview.WEBVIEW_HINT_NONE)
	w.BindExtensions()
	w.SetHTML(html)
	w.Run()
}
