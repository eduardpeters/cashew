package main

import (
	"bufio"
	"io"
	"log"
	"net"

	"github.com/eduardpeters/cashew/internal/commands"
	"github.com/eduardpeters/cashew/internal/store"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		log.Fatal("Error listening:", err)
	}
	defer listener.Close()

	log.Println("Now listening at port 6379")

	store := store.NewStore()

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
		}

		go handleConnection(store, conn)
	}
}

func handleConnection(store *store.Store, conn net.Conn) {
	defer conn.Close()
	addr := conn.RemoteAddr().String()
	log.Printf("accepted %s", addr)

	reader := bufio.NewReader(conn)

	for {
		args, err := commands.ParseCommand(reader)
		if err != nil {
			if err == io.EOF {
				log.Printf("%s disconnected", addr)
			} else {
				log.Printf("Read error from %s: %v", addr, err)
			}
			return
		}
		if len(args) == 0 {
			continue
		}

		result, err := commands.HandleCommand(store, args)
		if err != nil {
			result = commands.ResultError(err)
		}
		if _, err := conn.Write([]byte(result.Content)); err != nil {
			log.Printf("write error to %s: %v", addr, err)
			return
		}
	}
}
