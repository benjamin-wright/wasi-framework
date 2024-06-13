package main

import (
	"golang.org/x/sys/unix"
)

func main() {
	// js.Global().Set("hello", js.FuncOf(hello))
	unix.Exit(0)
}

// func hello(this js.Value, p []js.Value) interface{} {
// 	return "Hello from Go!"
// }
