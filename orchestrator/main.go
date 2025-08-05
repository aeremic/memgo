package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os/exec"
	"strconv"
	"strings"
)

// TODO: put in json config
var PATH string      // full path to exe
var HOST string      // host addr; should be changable also for each node
var NUM_OF_NODES int // TODO: min max
var INITIAL_PORT int // TODO: should be configurable for each node

const (
	STOP = "STOP"
)

func handleConnection(ctx context.Context, cancel context.CancelFunc, conn net.Conn) {
	defer conn.Close()

	reader := bufio.NewReader(conn)
	for {
		select {
		case <-ctx.Done():
			return
		default:
		}

		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				continue
			}
			fmt.Printf("%v", err)
			return
		}

		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			continue
		}

		command := strings.ToUpper(args[0])
		switch command {
		case STOP:
			for i := 0; i < NUM_OF_NODES; i++ {
				tcpAddr, err := net.ResolveTCPAddr("tcp", HOST+":"+strconv.Itoa(INITIAL_PORT+i))
				if err != nil {
					// TODO: Log
					log.Fatal(err)
					return
				}

				nodeconn, err := net.DialTCP("tcp", nil, tcpAddr)
				if err != nil {
					// TODO: Log
					log.Fatal(err)
					return
				}

				res, err := nodeconn.Write([]byte(STOP + "\n"))
				if err != nil {
					// TODO: Log
					log.Fatal(err)
					return
				}

				if res == 0 {
					// TODO: Log
					log.Fatalf("Result from node command is 0")
					return
				}
			}

			fmt.Printf("%s command received. Stopping thread and terminating orchestrator..\n", STOP)
			cancel()
			conn.Write([]byte(fmt.Sprintf("%s\n", "SUCCESS")))

			return
		default:
			// TODO: Find node to send the request
			return
		}
	}
}

func main() {
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

	fmt.Printf("Successfully started %d of nodes..\n", NUM_OF_NODES)

	// TODO: Stop processees
	// TODO: Run listener for orch

	ctx, cancel := context.WithCancel(context.Background())

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", HOST, strconv.Itoa(INITIAL_PORT-1)))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("Listening on ", listener.Addr())

	go func() {
		<-ctx.Done()
		fmt.Println("Shuting down hub..")
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return
			default:
				return
			}
		}

		go handleConnection(ctx, cancel, conn)
	}
}
