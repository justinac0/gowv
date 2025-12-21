package gowv_test

import (
	"fmt"
	"gowv/pkg/gowv"
	"testing"
)

func TestWebView(t *testing.T) {
	var w gowv.Instance
	w.Create(true, nil)
	w.Destroy()

	fmt.Printf("w: %v\n", w)
}
