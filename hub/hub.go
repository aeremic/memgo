package hub

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strings"
)

type ID string

const (
	STOP = "STOP"
	GET  = "GET"
	SET  = "SET"
)

type Command struct {
	Id ID
}

type Hub struct {
	host string
	port string
}

func New(host, port string) *Hub {
	return &Hub{
		host: host,
		port: port,
	}
}

func (h *Hub) Run() {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", h.host, h.port))
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	for {
		fmt.Printf("Serving client..\n")

		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		defer conn.Close()

		go func() {
			reader := bufio.NewReader(conn)
			for {
				command, err := reader.ReadString('\n')
				if err != nil {
					break
				}

				command = strings.TrimSpace(string(command))
				switch command {
				case STOP:
					fmt.Printf("%s command received. Stopping thread..\n", STOP)
					break
				case GET:
					fmt.Printf("%s command received.\n", GET)
				case SET:
					fmt.Printf("%s command received.\n", SET)
				default:
					log.Print("Unsupported command. Stopping thread..\n")
					break
				}

				fmt.Printf("incoming cmd: %s\n", command)
				conn.Write([]byte(fmt.Sprintf("following cmd received: %s\n", command)))
			}
		}()
	}
}
