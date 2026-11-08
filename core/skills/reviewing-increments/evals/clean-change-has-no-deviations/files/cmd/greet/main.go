package main

import (
	"fmt"
	"os"

	"example.com/greeter/internal/greet"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("greet: error args: missing-name")
		os.Exit(2)
	}
	fmt.Println("greet: " + greet.Line(os.Args[1]))
}
