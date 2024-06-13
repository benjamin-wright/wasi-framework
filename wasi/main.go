package main

import (
	"io"
	"os"
	"time"
)

func main() {
	input, err := io.ReadAll(os.Stdin)
	if err != nil {
		panic(err)
	}

	time.Sleep(5 * time.Second)

	output := []byte("Hello, " + string(input) + "!\n")

	_, err = os.Stdout.Write(output)
	if err != nil {
		panic(err)
	}
}
