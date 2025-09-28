package hub

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	common "memgo_common"
	"memgo_node/storage"
	"net"
	"strings"
)

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

func handleConnection(ctx context.Context, cancel context.CancelFunc, conn net.Conn,
	storage *storage.Storage) {
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

		line = strings.TrimSpace(line)
		if line == "" {
			log.Println("Empty command error")
			continue
		}

		args := strings.Fields(line)
		if len(args) == 0 {
			log.Println("Invalid number of arguments")
			continue
		}

		command := strings.ToUpper(args[0])
		switch command {
		case common.STOP:
			fmt.Printf("%s command received. Stopping thread and terminating hub..\n", common.STOP)
			cancel()
			conn.Write([]byte(fmt.Sprintf("%s\n", common.SUCCESS)))
			return
		case common.GETALL:
			fmt.Printf("%s command received.\n", common.GETALL)
			if len(args) != 1 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.GetAll()
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case common.GET:
			fmt.Printf("%s command received.\n", common.GET)
			if len(args) != 2 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.Get(args[1])
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case common.SET:
			fmt.Printf("%s command received.\n", common.SET)
			if len(args) != 3 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.Set(args[1], args[2])
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case common.DELETEALL:
			fmt.Printf("%s command received.\n", common.DELETE)
			if len(args) != 1 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.DeleteAll()
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case common.DELETE:
			fmt.Printf("%s command received.\n", common.GET)
			if len(args) != 2 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.Delete(args[1])
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case common.GETBYKEYANDPATH:
			fmt.Printf("%s command received.\n", common.GETBYKEYANDPATH)
			if len(args) != 3 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.GetByKeyAndPath(args[1], args[2])
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
		case common.SELECTBYPATH:
			fmt.Printf("%s command received.\n", common.SELECTBYPATH)
			if len(args) != 2 {
				log.Printf("Unsupported command %s received.\n", command)
				continue
			}

			res := storage.SelectByPath(args[1])
			conn.Write([]byte(fmt.Sprintf("%s\n", res)))
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

	fmt.Printf(common.SUCCESS + "\n")
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
				return nil
			default:
				return err
			}
		}

		go handleConnection(ctx, cancel, conn, storage)
	}
}
