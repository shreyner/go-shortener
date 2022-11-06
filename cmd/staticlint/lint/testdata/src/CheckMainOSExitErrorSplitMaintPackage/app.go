package main

import (
	"fmt"
	"os"
)

func main() {
	SomeFunction()
	fmt.Println("Hello world")
	os.Exit(1) // want "call os.Exit in main"
}
