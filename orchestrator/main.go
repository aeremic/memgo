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
	NOT_FOUND    = "NOT_FOUND"
	ERROR        = "ERROR"
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
				res, err := NODE_CONNECTIONS[i].Conn.Write([]byte(line))
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
			switch command {
			case GETALL, SELECTBYPATH:
				// TODO: Aggregate properly instead of returning {}{}{} results
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

					msg := strings.TrimSpace(string(bytes))
					if msg != NOT_FOUND {
						out.WriteString(msg)
					}
				}

				res := out.String()
				if res != "" {
					conn.Write([]byte(fmt.Sprintf("%s\n", res)))
				} else {
					conn.Write([]byte(fmt.Sprintf("%s\n", NOT_FOUND)))
				}

				continue
			case DELETEALL:
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

					msg := strings.TrimSpace(string(bytes))
					if msg != SUCCESS {
						conn.Write([]byte(fmt.Sprintf("%s\n", ERROR)))
						continue
					}
				}

				conn.Write([]byte(fmt.Sprintf("%s\n", SUCCESS)))

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
	config, err := Get("config.json")
	if err != nil {
		log.Fatal("Config file not loaded properly.", err)
	}

	if len(config.Nodes) == 0 {
		log.Fatal("Nodes not defined properly in the config file.")
	}

	for i := 0; i < len(config.Nodes); i++ {
		tcpAddr, err := net.ResolveTCPAddr("tcp", fmt.Sprintf("%s:%s", config.Nodes[i].Url, config.Nodes[i].Port))
		if err != nil {
			log.Fatal(err)
		}

		conn, err := net.DialTCP("tcp", nil, tcpAddr)
		if err != nil {
			log.Fatal(err)
		}

		NODE_CONNECTIONS = append(NODE_CONNECTIONS, Node{Conn: conn, Reader: bufio.NewReader(conn)})
	}

	HOST := config.Url
	PORT := config.Port

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
