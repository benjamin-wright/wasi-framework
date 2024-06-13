package main

import (
	"fmt"
	"syscall/js"
)

func main() {
	fmt.Println("Hello world!")
	js.Global().Set("hello", js.FuncOf(hello))
}

func hello(this js.Value, p []js.Value) interface{} {
	return "Hello from Go!"
}
