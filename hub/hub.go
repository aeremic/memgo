package hub

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"memgo/storage"
	"net"
	"strings"
)

type ID string

const (
	STOP   = "STOP"
	GET    = "GET"
	GETALL = "GETALL"
	SET    = "SET"
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

func handleConnection(ctx context.Context, cancel context.CancelFunc, conn net.Conn, storage *storage.Storage) {
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
		case GETALL:
			fmt.Printf("%s command received.\n", GET)
			if len(args) != 1 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.GetAll()
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case GET:
			fmt.Printf("%s command received.\n", GET)
			if len(args) != 2 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.Get(args[1])
			if res == "" {
				conn.Write([]byte(fmt.Sprintf("NotFound.\n")))
				continue
			}

			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case SET:
			fmt.Printf("%s command received.\n", SET)
			if len(args) != 3 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			storage.Set(args[1], args[2])
			conn.Write([]byte(fmt.Sprintf("Success\n")))
		default:
			log.Printf("Unsupported command %s received. Stopping thread..\n", command)
			return
		}
	}
}

func (h *Hub) Run(ctx context.Context, cancel context.CancelFunc, storage *storage.Storage) error {
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

		go handleConnection(ctx, cancel, conn, storage)
	}
}
