package main

import (
	"fmt"

	"go-sdk/sh"
)

func main() {
	value := sh.Prompt("first? ")
	fmt.Println("entered", value)

	value = sh.Promptf("%s? ", "second")
	fmt.Println("entered", value)
}
