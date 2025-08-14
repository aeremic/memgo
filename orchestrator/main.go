package main

import (
	"bufio"
	"context"
	"fmt"
	"hash/fnv"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

type Node struct {
	Conn   *net.TCPConn
	Reader *bufio.Reader
}

var NODE_CONNECTIONS []Node

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
			fmt.Printf("Error on reading string: %v", err)
			return
		}

		preparedLine := strings.TrimSpace(line)
		if preparedLine == "" {
			log.Println("Empty command error")
			continue
		}

		args := strings.Fields(preparedLine)
		if len(args) == 0 {
			log.Println("Invalid number of arguments")
			continue
		}

		command := strings.ToUpper(args[0])
		if command == STOP {
			for i := 0; i < len(NODE_CONNECTIONS); i++ {
				res, err := NODE_CONNECTIONS[i].Conn.Write([]byte(STOP + "\n"))
				if err != nil {
					log.Println(err)
					return
				}

				if res == 0 {
					log.Println("Result from node command is 0")
					return
				}
			}

			fmt.Printf("%s command received. Stopping thread and terminating orchestrator..\n", STOP)
			cancel()
			conn.Write([]byte(fmt.Sprintf("%s\n", SUCCESS)))

			return
		} else {
			// TODO: Issue for getall and commands that don't have keys; GETALL, DELETEALL and SELECTBYPATH
			// should agreggate results

			if len(args) < 2 {
				log.Println("Invalid number of arguments")
				continue
			}

			// Consistent hashing
			key := args[1]
			h := fnv.New32a()
			h.Write([]byte(key))

			node_index := int(h.Sum32()) % len(NODE_CONNECTIONS)
			res, err := NODE_CONNECTIONS[node_index].Conn.Write([]byte(line))
			if err != nil {
				log.Println(err)
				continue
			}

			if res == 0 {
				log.Println("Result from node command is 0")
				continue
			}

			bytes, err := NODE_CONNECTIONS[node_index].Reader.ReadBytes('\n')
			if err != nil {
				log.Println("Result from node command is invalid")
				continue
			}

			conn.Write(bytes)

			continue
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

		NODE_CONNECTIONS = append(NODE_CONNECTIONS, Node{Conn: conn, Reader: bufio.NewReader(conn)})
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
