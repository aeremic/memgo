package hub

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
			fmt.Printf("%s command received. Stopping thread and terminating hub..\n", STOP)
			cancel()
			return
		case GET:
			fmt.Printf("%s command received.\n", GET)
		case SET:
			fmt.Printf("%s command received.\n", SET)
		default:
			log.Printf("Unsupported command %s received. Stopping thread..\n", command)
			return
		}

		// fmt.Printf("incoming cmd: %s\n", command)
		// conn.Write([]byte(fmt.Sprintf("following cmd received: %s\n", command)))
	}
}

func (h *Hub) Run(ctx context.Context, cancel context.CancelFunc) error {
	listener, err := net.Listen("tcp", fmt.Sprintf("%s:%s", h.host, h.port))
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
		fmt.Printf("Serving client..\n")

		conn, err := listener.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				return nil
			default:
				return err
			}
		}

		go handleConnection(ctx, cancel, conn)
	}
}
