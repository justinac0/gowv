package main

import (
	webview "gowv"
	"time"
)

const html = `
<button id="set_fullscreen">Set Fullscreen</button>
<button id="set_icon">Set Icon</button>
<button id="show">Show</button>
<button id="hide">Hide</button>
<button id="set_maximized">Set Maximized</button>
<button id="set_minimized">Set Minimized</button>
<script>
  const set_fullscreen = document.querySelector("#set_fullscreen")
  const set_icon = document.querySelector("#set_icon")
  const show = document.querySelector("#show")
  const hide = document.querySelector("#hide")
  const set_maximized = document.querySelector("#set_maximized")
  const set_minimized = document.querySelector("#set_minimized")
  document.addEventListener("DOMContentLoaded", () => {
  	set_fullscreen.addEventListener("click", () => { window.native_set_fullscreen().then(result => {}); });
    set_icon.addEventListener("click", () => { window.native_set_icon().then(result => {}); });
    show.addEventListener("click", () => { window.native_show().then(result => {}); });
    hide.addEventListener("click", () => { window.native_hide().then(result => {}); });
    set_maximized.addEventListener("click", () => { window.native_set_maximized().then(result => {}); });
    set_minimized.addEventListener("click", () => { window.native_set_minimized().then(result => {}); });
  });
</script>`

func main() {
	w := webview.Instance{}
	w.Create(true, nil)
	defer w.Destroy()
	w.SetTitle("Extension Example")
	w.SetSize(480, 320, webview.WEBVIEW_HINT_NONE)

	exTimeout := 2 * time.Second

	w.Bind("native_set_fullscreen", func() {
		go func() {
			w.SetFullscreen(true)
			time.Sleep(exTimeout)
			w.SetFullscreen(false)
		}()
	})

	w.Bind("native_set_icon", func() {
		w.SetIcon("./examples/extension/icon.ico")
	})

	w.Bind("native_show", func() {
		w.Show()
	})

	w.Bind("native_hide", func() {
		go func() {
			w.Hide()
			time.Sleep(exTimeout)
			w.Show()
		}()
	})

	w.Bind("native_set_maximized", func() {
		w.SetMaximized()
	})

	w.Bind("native_set_minimized", func() {
		w.SetMinimized()
	})

	w.SetHTML(html)
	w.Run()
}
