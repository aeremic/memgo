package main

import (
	"bufio"
	"bytes"
	"context"
	"errors"
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

// TODO: Make common package for command enums
const (
	STOP         = "STOP"
	SUCCESS      = "SUCCESS"
	GETALL       = "GETALL"
	DELETEALL    = "DELETEALL"
	SELECTBYPATH = "SELECTBYPATH"
)

func getNodeIndex(key string) int {
	h := fnv.New32a()
	h.Write([]byte(key))

	return int(h.Sum32()) % len(NODE_CONNECTIONS)
}

func forwardMsg(node_index int, line string) error {
	res, err := NODE_CONNECTIONS[node_index].Conn.Write([]byte(line))
	if err != nil {
		return err
	}

	if res == 0 {
		return errors.New("Result from node command is 0")
	}

	return nil
}

func receiveMsg(node_index int) ([]byte, error) {
	bytes, err := NODE_CONNECTIONS[node_index].Reader.ReadBytes('\n')
	if err != nil {
		return []byte{}, err
	}

	return bytes, nil
}

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
			// TODO: Issue for getall and commands that don't have keys; GETALL, DELETEALL and SELECTBYPATH should agreggate results
			// TODO: getall {"aaa":"bbb"}{"test":"aaa"} and getall NOTFOUND {"test":"aaa"} fix
			switch command {
			case GETALL, DELETEALL, SELECTBYPATH:
				var out bytes.Buffer
				for i := 0; i < len(NODE_CONNECTIONS); i++ {
					err := forwardMsg(i, line)
					if err != nil {
						log.Printf("Error on forwarding message: %v", err)
						continue
					}

					bytes, err := receiveMsg(i)
					if err != nil {
						log.Printf("Error on receiving message: %v", err)
						continue
					}

					out.WriteString(string(bytes))
				}

				res := out.String()
				if res != "" {
					conn.Write([]byte(fmt.Sprintf("%s", res)))
				}

				continue
			default:
				if len(args) < 2 {
					log.Println("Invalid number of arguments")
					continue
				}

				// Consistent hashing
				// TODO: Investigate adding/removing node solutions
				node_index := getNodeIndex(args[1])
				err := forwardMsg(node_index, line)
				if err != nil {
					log.Printf("Error on forwarding message: %v", err)
					continue
				}

				bytes, err := receiveMsg(node_index)
				if err != nil {
					log.Printf("Error on receiving message: %v", err)
					continue
				}

				conn.Write(bytes)

				continue
			}
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
				log.Println(err)
				return
			}
		}

		go handleConnection(ctx, cancel, conn)
	}
}
