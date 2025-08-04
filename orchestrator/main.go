package main

import (
	"fmt"
	"log"
	"os/exec"
	"strconv"
)

func main() {
	// TODO: put in json config
	var PATH string      // full path to exe
	var HOST string      // host addr; should be changable also for each node
	var NUM_OF_NODES int // TODO: min max
	var INITIAL_PORT int // TODO: should be configurable for each node

	PATH = "/Users/eremic/memgo/node/memgo_node"
	HOST = "localhost"
	NUM_OF_NODES = 2
	INITIAL_PORT = 1234

	for i := 0; i < NUM_OF_NODES; i++ {
		cmd := exec.Command(PATH, HOST, strconv.Itoa(INITIAL_PORT+i))

		if err := cmd.Start(); err != nil {
			log.Fatal(err)
			break
		}
	}

	fmt.Printf("Successfully started %d of nodes..", NUM_OF_NODES)
}
