package main

import (
	"fmt"

	"golang.org/x/example/stringutil"
)

func main() {
	original := "Hello, OTUS!"
	reversed := stringutil.Reverse(original)
	fmt.Println(reversed)
}
