package main

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var NODE_CONNECTIONS []*net.TCPConn

const (
	STOP    = "STOP"
	SUCCESS = "SUCCESS"
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
			for i := 0; i < len(NODE_CONNECTIONS); i++ {
				res, err := NODE_CONNECTIONS[i].Write([]byte(STOP + "\n"))
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
			conn.Write([]byte(fmt.Sprintf("%s\n", SUCCESS)))

			return
		default:
			// TODO: Find node to send the request
			return
		}
	}
}

func main() {
	nodes := [2]string{"localhost:1234", "localhost:1235"} // TODO: Pull from config file

	for i := 0; i < len(nodes); i++ {
		tcpAddr, err := net.ResolveTCPAddr("tcp", nodes[i])
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatal(err)
		}

		NODE_CONNECTIONS = append(NODE_CONNECTIONS, conn)
	}

	var PORT string
	var HOST string

	args := os.Args
	if len(args) == 3 {
		HOST = args[1]
		PORT = args[2]
	} else {
		fmt.Printf("Port number and host are not provided, default ones will be used (localhost:1233).\n")
		HOST = "localhost"
		PORT = "1233"
	}

	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", HOST, PORT))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Printf(SUCCESS + "\n")
	fmt.Println("Listening on ", listener.Addr())

	ctx, cancel := context.WithCancel(context.Background())

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
