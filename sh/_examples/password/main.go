package main

import (
	"fmt"

	"go-sdk/sh"
)

func main() {
	value, err := sh.Password("first? ")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)

	value, err = sh.Passwordf("%s? ", "second")
	if err != nil {
		sh.Fatal(err)
	}
	fmt.Println("entered", value)
}
