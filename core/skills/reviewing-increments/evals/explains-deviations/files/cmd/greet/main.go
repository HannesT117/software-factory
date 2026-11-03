package main

import (
	"fmt"
	"os"
	"sort"
	"strings"

	"example.com/greeter/internal/greet"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("greet: error args: missing-name")
		os.Exit(2)
	}
	if os.Args[1] == "--known" {
		data, err := os.ReadFile("names.txt")
		if err != nil {
			fmt.Println("greet: error names.txt: " + err.Error())
			os.Exit(1)
		}
		names := strings.Fields(string(data))
		sort.Strings(names)
		for _, n := range names {
			fmt.Println("greet: " + n + " is known")
		}
		return
	}
	fmt.Println("greet: " + greet.Line(os.Args[1]))
}
