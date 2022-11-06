package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hello world")
	if true {
		os.Exit(1) // want "call os.Exit in main"
	}
}
