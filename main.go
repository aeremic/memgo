package main

import (
	"context"
	"fmt"
	"log"
	"memgo/hub"
	"memgo/storage"
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
		fmt.Printf("Port number and host are not provided, default ones will be used (localhost:1234).\n")
		HOST = "localhost"
		PORT = "1234"
	}

	fmt.Printf("Starting hub on: %s:%s\n", HOST, PORT)

	s := storage.New()
	ctx, cancel := context.WithCancel(context.Background())

	h := hub.New(HOST, PORT)
	if err := h.Run(ctx, cancel, s); err != nil {
		log.Fatal(err)
	}
}
