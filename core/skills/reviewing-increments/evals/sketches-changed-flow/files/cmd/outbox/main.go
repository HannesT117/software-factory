package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 3 || os.Args[1] != "push" {
		fmt.Println("outbox: error args: usage outbox push <url>")
		os.Exit(2)
	}
	url := os.Args[2]
	data, err := os.ReadFile("outbox.txt")
	if os.IsNotExist(err) {
		fmt.Println("outbox: sent 0")
		return
	}
	lines := strings.Split(strings.TrimRight(string(data), "\n"), "\n")
	for _, msg := range lines {
		status := send(url, msg)
		if status != http.StatusOK {
			// The server may be busy for a moment: wait and try once more.
			time.Sleep(time.Second)
			status = send(url, msg)
		}
		if status != http.StatusOK {
			fmt.Printf("outbox: error send: %d\n", status)
			os.Exit(1)
		}
	}
	os.Remove("outbox.txt")
	fmt.Printf("outbox: sent %d\n", len(lines))
}

func send(url, msg string) int {
	resp, err := http.Post(url, "text/plain", strings.NewReader(msg))
	if err != nil {
		return 0
	}
	resp.Body.Close()
	return resp.StatusCode
}
