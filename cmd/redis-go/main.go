package main

import (
	"fmt"
	"log"
	"net"
	"strings"

	"github.com/codecrafters-io/redis-starter-go/internal/db"
	"github.com/codecrafters-io/redis-starter-go/internal/resp"
)

var store = db.New()

func main() {
	log.SetFlags(0)

	listener, err := net.Listen("tcp", "0.0.0.0:6379") //nolint:gosec // G102: must listen on all interfaces for the tester
	if err != nil {
		log.Fatalf("failed to bind to port 6379: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("error accepting connection: %v", err)
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer func() {
		if err := conn.Close(); err != nil {
			log.Printf("error closing connection: %v", err)
		}
	}()

	buf := make([]byte, 1024)
	for {
		n, err := conn.Read(buf)
		if err != nil {
			log.Printf("error reading from connection: %v", err)
			break
		}

		args, err := resp.Parse(buf[:n])
		if err != nil || len(args) == 0 {
			log.Printf("error parsing input: %v", err)
			return
		}

		switch strings.ToUpper(args[0]) {
		case "PING":
			err = writeToConn(conn, resp.SimpleString("PONG"))

		case "ECHO":
			if len(args) < 2 {
				log.Printf("ECHO call needs %d args but got %d args", 2, len(args))
				continue
			}

			err = writeToConn(conn, resp.BulkString(args[1]))

		case "SET":
			if len(args) < 3 {
				log.Printf("SET call needs %d args but got %d args", 3, len(args))
				continue
			}

			store.Set(args[1], args[2])
			err = writeToConn(conn, resp.SimpleString("OK"))

		case "GET":
			if len(args) < 2 {
				log.Printf("GET call needs %d args but got %d args", 2, len(args))
				continue
			}

			value := store.Get(args[1])
			err = writeToConn(conn, resp.BulkString(value))
		}

		if err != nil {
			log.Print(err)
			break
		}
	}
}

func writeToConn(conn net.Conn, data string) error {
	_, err := conn.Write([]byte(data))
	if err != nil {
		return fmt.Errorf("error writing a response: %w", err)
	}

	return nil
}
