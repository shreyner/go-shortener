package random

import "fmt"

func ExampleRandSeq() {
	length := 5
	s := RandSeq(length)

	fmt.Println(s) // O7Yl7
}
