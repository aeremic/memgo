package main

import (
	"fmt"
	"memgo/hub"
	"os"
)

func main() {
	var PORT string
	var HOST string

	args := os.Args
	if len(args) == 3 {
		HOST = args[1]
		PORT = args[2]
	} else {
		fmt.Printf("Port number and host are not provided, default ones will be used (localhost:1234).")
		HOST = "localhost"
		PORT = "1234"
	}

	fmt.Printf("Listening on: %s:%s", HOST, PORT)

	h := hub.New(HOST, PORT)
	h.Run()
}
