// package main

// import (
// 	"gowv"
// )

// type Counter struct {
// 	Count uint `json:"count"`
// }

// func main() {
// 	w := gowv.Instance{}

// 	w.Create(true, nil)

// 	defer func() {
// 		gowv.PanicOnError(w.Destroy())
// 	}()

// 	gowv.PanicOnError(w.SetTitle("Bind Example"))
// 	gowv.PanicOnError(w.SetSize(480, 320, gowv.WEBVIEW_HINT_NONE))

// 	html := `
// 	<div>
//   		<button id="increment">+</button>
//     	<button id="decrement">−</button>
//      	<span>Counter: <span id="counterResult">0</span></span>
//     </div>
//     <hr />
//     <div>
//     	<button id="compute">Compute</button>
//      	<span>Result: <span id="computeResult">(not started)</span></span>
//     </div>
//     <script type=\"module\">
//     	const getElements = ids => Object.assign({}, ...ids.map(
//      		id => ({ [id]: document.getElementById(id) })));
//       		const ui = getElements([
//         		"increment", "decrement", "counterResult", "compute",
//           		"computeResult"
//             ]);
//             ui.increment.addEventListener("click", async () => {
//             	ui.counterResult.textContent = await window.count(1);
//             });
//             ui.decrement.addEventListener("click", async () => {
//              	ui.counterResult.textContent = await window.count(-1);
//             });
//             ui.compute.addEventListener("click", async () => {
//             	ui.compute.disabled = true;
//              	ui.computeResult.textContent = "(pending)";
//               	ui.computeResult.textContent = await window.compute(6, 7);
//                	ui.compute.disabled = false;
//             });
//     </script>`

// 	var count uint = 0

// 	gowv.PanicOnError(w.Bind("count", func() Counter {
// 		count++
// 		return Counter{Count: count}
// 	}))
// 	gowv.PanicOnError(w.Bind("compute", func() Counter {
// 		return Counter{Count: count}
// 	}))

//		gowv.PanicOnError(w.SetHTML(html))
//		gowv.PanicOnError(w.Run())
//	}
package main

import webview "gowv"

const html = `<button id="increment">Tap me</button>
<div>You tapped <span id="count">0</span> time(s).</div>
<script>
  const [incrementElement, countElement] =
    document.querySelectorAll("#increment, #count");
  document.addEventListener("DOMContentLoaded", () => {
    incrementElement.addEventListener("click", () => {
      window.increment().then(result => {
        countElement.textContent = result.count;
      });
    });
  });
</script>`

type IncrementResult struct {
	Count uint `json:"count"`
}

func main() {
	var count uint = 0
	w := webview.Instance{}
	w.Create(true, nil)
	defer w.Destroy()
	w.SetTitle("Bind Example")
	w.SetSize(480, 320, webview.WEBVIEW_HINT_NONE)

	// A binding that increments a value and immediately returns the new value.
	w.Bind("increment", func() IncrementResult {
		count++
		return IncrementResult{Count: count}
	})

	w.SetHTML(html)
	w.Run()
}
