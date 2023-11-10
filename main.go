package main

import (
	"fmt"
)

func main() {
	filename := "target/computer.go"
	needle := `TestString`

	ret := checkForStringInFile(filename, needle)
	fmt.Printf("ret: %v\n", ret)
}
